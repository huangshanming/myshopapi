package main

// @title           mymall 商品服务 API
// @version         1.0
// @description     商品与分类查询
// @host            localhost:9080
// @BasePath        /

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"mymall/pkg/cache"
	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	"mymall/pkg/jwt"
	applog "mymall/pkg/log"
	"mymall/pkg/metrics"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/pkg/mq"
	"mymall/pkg/telemetry"
	contenthandler "mymall/services/catalog-service/internal/content/handler"
	contentlogic "mymall/services/catalog-service/internal/content/logic"
	contentmodel "mymall/services/catalog-service/internal/content/model"
	catalogmq "mymall/services/catalog-service/internal/mq"
	notifyhandler "mymall/services/catalog-service/internal/notify/handler"
	notifymodel "mymall/services/catalog-service/internal/notify/model"
	producthandler "mymall/services/catalog-service/internal/product/handler"
	productlogic "mymall/services/catalog-service/internal/product/logic"
	productmodel "mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/server"
	shopopshandler "mymall/services/catalog-service/internal/shopops/handler"
	shopopsmodel "mymall/services/catalog-service/internal/shopops/model"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/uploadpath"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/catalog-service.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}

	logger, err := applog.New("catalog-service")
	if err != nil {
		log.Fatalf("初始化日志失败：%v", err)
	}
	defer logger.Sync()

	ctx := context.Background()
	shutdownTrace, err := telemetry.Init(ctx, cfg.Telemetry)
	if err != nil {
		logger.Warn("telemetry init skipped")
	}
	defer shutdownTrace(context.Background())

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	if err := database.AutoMigrateIfDebug(cfg.Server.Mode, db,
		&productmodel.Product{},
		&productmodel.ProductCategory{},
		&productmodel.ProductSku{},
		&productmodel.ProductImage{},
		&productmodel.ProductTag{},
		&productmodel.ProductTagRel{},
		&productmodel.ProductAttrTemplate{},
		&productmodel.ProductAttr{},
		&productmodel.ProductSchedule{},
		&productmodel.ProductBatchJob{},
		&productmodel.ProductOpLog{},
		&productmodel.ProductFavorite{},
		&shopopsmodel.ShopRole{},
		&shopopsmodel.ShopMenu{},
		&shopopsmodel.ShopRoleMenu{},
		&shopopsmodel.ShopUserRole{},
		&contentmodel.CommunityArticle{},
		&contentmodel.CommunityArticleCategory{},
		&contentmodel.CommunityArticleComment{},
		&contentmodel.CommunityCommentEmoji{},
		&contentmodel.CommunityArticleImg{},
		&contentmodel.ArticleLike{},
		&contentmodel.ArticleFavorite{},
		&contentmodel.ArticleAudience{},
		&contentmodel.HomepageBanner{},
		&notifymodel.ShopNotification{},
	); err != nil {
		log.Fatalf("AutoMigrate 失败：%v", err)
	}

	var redisClient *redis.Client
	rc, err := cache.NewRedis(cfg.Redis)
	if err != nil {
		logger.Warn("redis unavailable, cache disabled")
	} else {
		redisClient = rc
	}

	var mqClient *mq.Client
	mqc, err := mq.New(cfg.RabbitMQ)
	if err != nil {
		logger.Warn("rabbitmq unavailable")
	} else {
		mqClient = mqc
		defer mqc.Close()
	}

	svcCtx := svc.NewServiceContext(cfg, db, redisClient, mqClient)
	if err := svcCtx.ShopRBAC.EnsureShopMenus(); err != nil {
		logger.Warn(fmt.Sprintf("seed shop menus: %v", err))
	} else {
		logger.Info("shop menus seeded (layered)")
	}
	catalogLogic := productlogic.NewCatalogLogic(context.Background(), svcCtx)
	catalogPublic := producthandler.NewCatalogPublicHandler(svcCtx)
	catalogAdmin := producthandler.NewCatalogAdminHandler(svcCtx)
	favoriteUser := producthandler.NewFavoriteUserHandler(svcCtx)
	favoriteAdmin := producthandler.NewFavoriteAdminHandler(svcCtx)
	adminH := producthandler.NewProductAdminHandler(svcCtx)
	shopOpsH := shopopshandler.NewShopOpsHandler(svcCtx)
	articleAdminH := contenthandler.NewArticleAdminHandler(svcCtx)
	articleMerchantH := contenthandler.NewArticleMerchantHandler(svcCtx)
	articlePublicH := contenthandler.NewArticlePublicHandler(svcCtx)
	shopUploadH := producthandler.NewShopUploadHandler()
	platformProductH := producthandler.NewPlatformProductHandler(svcCtx)
	notifH := notifyhandler.NewNotificationHandler(svcCtx)
	productAdminLogic := productlogic.NewProductAdminLogic(context.Background(), svcCtx)
	articleLogic := contentlogic.NewArticleLogic(context.Background(), svcCtx)

	// 商品定时上下架 + 文章定时发布（同进程 Mutex 防叠跑）
	go func() {
		var scheduleMu sync.Mutex
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for range t.C {
			scheduleMu.Lock()
			productAdminLogic.RunSchedules()
			articleLogic.RunPublishSchedules()
			scheduleMu.Unlock()
		}
	}()

	if svcCtx.MQ != nil {
		consumer := catalogmq.NewConsumer(svcCtx.MQ, catalogLogic, logger)
		if err := consumer.Start(); err != nil {
			logger.Warn("mq consumer start failed")
		}
	}

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(ctx)
	})
	if redisClient != nil {
		healthReg.Register("redis", func(ctx context.Context) error {
			return cache.Ping(ctx, redisClient)
		})
	}
	if mqClient != nil {
		healthReg.Register("rabbitmq", mqClient.Ping)
	}

	rpcServer := server.StartZRPC(cfg.Server.GRPCPort, catalogLogic, logger)
	go func() {
		logger.Info(fmt.Sprintf("catalog-service zRPC 启动 :%d", cfg.Server.GRPCPort))
		rpcServer.Start()
	}()
	defer rpcServer.Stop()

	serverHTTP := httpserver.NewRest(cfg.Server.HTTPPort, cfg.Server.Mode)
	defer serverHTTP.Stop()

	rid := middleware.RequestID()
	gwShop := middleware.GatewayIdentity(true)
	gwUser := middleware.GatewayIdentity(false)
	gw := middleware.GatewayIdentity(false)
	merchantRoles := middleware.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)
	adminRoles := middleware.RequireRoles(jwt.RolePlatformAdmin)

	serverHTTP.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: rid(httpserver.Healthz("catalog-service"))},
		{Method: http.MethodGet, Path: "/readyz", Handler: rid(healthReg.ReadyHandler())},
		{Method: http.MethodGet, Path: "/metrics", Handler: rid(metrics.Handler())},

		{Method: http.MethodGet, Path: "/api/v1/products/list", Handler: rid(catalogPublic.GetProductList)},
		{Method: http.MethodGet, Path: "/api/v1/products/detail", Handler: rid(catalogPublic.GetProductDetail)},
		{Method: http.MethodGet, Path: "/api/v1/products/sales-rank", Handler: rid(catalogPublic.GetSalesRank)},
		{Method: http.MethodGet, Path: "/api/v1/product_category/list", Handler: rid(catalogPublic.GetCategoryList)},
		{Method: http.MethodGet, Path: "/api/v1/product_category/detail", Handler: rid(catalogPublic.GetCategoryDetail)},

		{Method: http.MethodGet, Path: "/api/v1/banners", Handler: rid(articlePublicH.ListBanners)},

		{Method: http.MethodGet, Path: "/api/v1/articles/list", Handler: rid(articlePublicH.List)},
		{Method: http.MethodGet, Path: "/api/v1/articles/:id", Handler: rid(articlePublicH.Detail)},
		{Method: http.MethodGet, Path: "/api/v1/articles/:id/comments", Handler: rid(articlePublicH.ListComments)},
		{Method: http.MethodPost, Path: "/api/v1/articles/:id/comments", Handler: rid(gw(articlePublicH.CreateComment))},
		{Method: http.MethodGet, Path: "/api/v1/comment-emojis", Handler: rid(articlePublicH.ListEmojis)},
		{Method: http.MethodPost, Path: "/api/v1/articles/:id/like", Handler: rid(gw(articlePublicH.Like))},
		{Method: http.MethodDelete, Path: "/api/v1/articles/:id/like", Handler: rid(gw(articlePublicH.Unlike))},
		{Method: http.MethodPost, Path: "/api/v1/articles/:id/favorite", Handler: rid(gw(articlePublicH.Favorite))},
		{Method: http.MethodDelete, Path: "/api/v1/articles/:id/favorite", Handler: rid(gw(articlePublicH.Unfavorite))},
		{Method: http.MethodGet, Path: "/api/v1/articles/:id/engagement", Handler: rid(gw(articlePublicH.Status))},
		{Method: http.MethodGet, Path: "/api/v1/user/article-favorites", Handler: rid(gw(articlePublicH.ListMyFavorites))},
		{Method: http.MethodGet, Path: "/api/v1/user/article-likes", Handler: rid(gw(articlePublicH.ListMyLikes))},
		{Method: http.MethodGet, Path: "/api/v1/user/articles", Handler: rid(gw(articlePublicH.ListMine))},
		{Method: http.MethodPost, Path: "/api/v1/user/articles", Handler: rid(gw(articlePublicH.CreateMine))},
		{Method: http.MethodGet, Path: "/api/v1/user/articles/:id", Handler: rid(gw(articlePublicH.DetailMine))},
		{Method: http.MethodPut, Path: "/api/v1/user/articles/:id", Handler: rid(gw(articlePublicH.UpdateMine))},
		{Method: http.MethodDelete, Path: "/api/v1/user/articles/:id", Handler: rid(gw(articlePublicH.DeleteMine))},
		{Method: http.MethodPost, Path: "/api/v1/user/article-uploads", Handler: rid(gw(articlePublicH.UploadMine))},

		{Method: http.MethodPost, Path: "/api/v1/user/favorites", Handler: rid(gw(favoriteUser.Add))},
		{Method: http.MethodDelete, Path: "/api/v1/user/favorites/:product_id", Handler: rid(gw(favoriteUser.Remove))},
		{Method: http.MethodPost, Path: "/api/v1/user/favorites/batch-remove", Handler: rid(gw(favoriteUser.RemoveBatch))},
		{Method: http.MethodGet, Path: "/api/v1/user/favorites", Handler: rid(gw(favoriteUser.List))},
		{Method: http.MethodGet, Path: "/api/v1/products/:id/favorite", Handler: rid(gw(favoriteUser.Status))},
		{Method: http.MethodGet, Path: "/api/v1/products/:id/favorite-count", Handler: rid(favoriteUser.Count)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id/favorites", Handler: rid(middleware.Chain(favoriteAdmin.AdminUserList, gwUser, adminRoles))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/products", Handler: rid(middleware.Chain(adminH.List, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products", Handler: rid(middleware.Chain(adminH.Create, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/batch", Handler: rid(middleware.Chain(adminH.Batch, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/products/jobs/:id", Handler: rid(middleware.Chain(adminH.JobStatus, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/recycle/restore", Handler: rid(middleware.Chain(adminH.RecycleRestore, gwShop, merchantRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/products/recycle", Handler: rid(middleware.Chain(adminH.RecycleDelete, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/products/export", Handler: rid(middleware.Chain(adminH.Export, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/import", Handler: rid(middleware.Chain(adminH.Import, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/products/op-logs", Handler: rid(middleware.Chain(adminH.OpLogs, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/products/:id", Handler: rid(middleware.Chain(adminH.Detail, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/products/:id", Handler: rid(middleware.Chain(adminH.Update, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/products/:id/status", Handler: rid(middleware.Chain(adminH.SetStatus, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/:id/copy", Handler: rid(middleware.Chain(adminH.Copy, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/:id/schedules", Handler: rid(middleware.Chain(adminH.Schedule, gwShop, merchantRoles))},

		{Method: http.MethodPut, Path: "/api/v1/merchant/skus/:id/stock", Handler: rid(middleware.Chain(adminH.AdjustStock, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/skus/batch-stock", Handler: rid(middleware.Chain(adminH.BatchStock, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/stocks/warnings", Handler: rid(middleware.Chain(adminH.StockWarnings, gwShop, merchantRoles))},

		{Method: http.MethodPost, Path: "/api/v1/merchant/uploads/images", Handler: rid(middleware.Chain(adminH.Upload, gwShop, merchantRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/schedules/:id", Handler: rid(middleware.Chain(adminH.CancelSchedule, gwShop, merchantRoles))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/tags", Handler: rid(middleware.Chain(adminH.ListTags, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/tags", Handler: rid(middleware.Chain(adminH.SaveTag, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/tags/:id", Handler: rid(middleware.Chain(adminH.SaveTag, gwShop, merchantRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/tags/:id", Handler: rid(middleware.Chain(adminH.DeleteTag, gwShop, merchantRoles))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/attr-templates", Handler: rid(middleware.Chain(adminH.ListAttrTemplates, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/attr-templates", Handler: rid(middleware.Chain(adminH.SaveAttrTemplate, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/attr-templates/:id", Handler: rid(middleware.Chain(adminH.SaveAttrTemplate, gwShop, merchantRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/attr-templates/:id", Handler: rid(middleware.Chain(adminH.DeleteAttrTemplate, gwShop, merchantRoles))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/auth/me", Handler: rid(middleware.Chain(shopOpsH.AuthMe, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shop/roles", Handler: rid(middleware.Chain(shopOpsH.ListRoles, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shop/roles/:id/menus", Handler: rid(middleware.Chain(shopOpsH.RoleMenus, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/shop/roles", Handler: rid(middleware.Chain(shopOpsH.SaveRole, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/shop/roles/:id", Handler: rid(middleware.Chain(shopOpsH.SaveRole, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shop/menus", Handler: rid(middleware.Chain(shopOpsH.ListMenus, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shop/staff", Handler: rid(middleware.Chain(shopOpsH.ListStaff, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/shop/staff", Handler: rid(middleware.Chain(shopOpsH.BindStaff, gwShop, merchantRoles))},

		{Method: http.MethodGet, Path: "/api/v1/admin/products", Handler: rid(middleware.Chain(platformProductH.List, gwUser, adminRoles))},
		{Method: http.MethodPut, Path: "/api/v1/admin/products/:id/off_sale", Handler: rid(middleware.Chain(platformProductH.OffSale, gwUser, adminRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/products/:id", Handler: rid(middleware.Chain(platformProductH.Delete, gwUser, adminRoles))},
		{Method: http.MethodGet, Path: "/api/v1/admin/categories", Handler: rid(middleware.Chain(catalogAdmin.AdminListCategories, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/categories", Handler: rid(middleware.Chain(catalogAdmin.AdminCreateCategory, gwUser, adminRoles))},
		{Method: http.MethodPut, Path: "/api/v1/admin/categories/:id", Handler: rid(middleware.Chain(catalogAdmin.AdminUpdateCategory, gwUser, adminRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/categories/:id", Handler: rid(middleware.Chain(catalogAdmin.AdminDeleteCategory, gwUser, adminRoles))},

		// 社区文章 — 平台
		{Method: http.MethodGet, Path: "/api/v1/admin/articles/stats", Handler: rid(middleware.Chain(articleAdminH.Stats, gwUser, adminRoles))},
		{Method: http.MethodGet, Path: "/api/v1/admin/articles/recycle", Handler: rid(middleware.Chain(articleAdminH.List, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/recycle/restore", Handler: rid(middleware.Chain(articleAdminH.RecycleRestore, gwUser, adminRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/articles/recycle", Handler: rid(middleware.Chain(articleAdminH.RecycleDelete, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/batch-audit", Handler: rid(middleware.Chain(articleAdminH.BatchAudit, gwUser, adminRoles))},
		{Method: http.MethodGet, Path: "/api/v1/admin/articles", Handler: rid(middleware.Chain(articleAdminH.List, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles", Handler: rid(middleware.Chain(articleAdminH.Create, gwUser, adminRoles))},
		{Method: http.MethodGet, Path: "/api/v1/admin/articles/:id", Handler: rid(middleware.Chain(articleAdminH.Detail, gwUser, adminRoles))},
		{Method: http.MethodPut, Path: "/api/v1/admin/articles/:id", Handler: rid(middleware.Chain(articleAdminH.Update, gwUser, adminRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/articles/:id", Handler: rid(middleware.Chain(articleAdminH.SoftDelete, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/:id/audit", Handler: rid(middleware.Chain(articleAdminH.Audit, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/:id/top", Handler: rid(middleware.Chain(articleAdminH.Top, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/:id/offline", Handler: rid(middleware.Chain(articleAdminH.Offline, gwUser, adminRoles))},
		{Method: http.MethodGet, Path: "/api/v1/admin/article-categories", Handler: rid(middleware.Chain(articleAdminH.CategoryList, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/article-categories", Handler: rid(middleware.Chain(articleAdminH.CategoryCreate, gwUser, adminRoles))},
		{Method: http.MethodPut, Path: "/api/v1/admin/article-categories/:id", Handler: rid(middleware.Chain(articleAdminH.CategoryUpdate, gwUser, adminRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/article-categories/:id", Handler: rid(middleware.Chain(articleAdminH.CategoryDelete, gwUser, adminRoles))},
		{Method: http.MethodGet, Path: "/api/v1/admin/article-comments", Handler: rid(middleware.Chain(articleAdminH.CommentList, gwUser, adminRoles))},
		{Method: http.MethodPatch, Path: "/api/v1/admin/article-comments/:id", Handler: rid(middleware.Chain(articleAdminH.CommentPatch, gwUser, adminRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/article-comments/:id", Handler: rid(middleware.Chain(articleAdminH.CommentDelete, gwUser, adminRoles))},
		{Method: http.MethodGet, Path: "/api/v1/admin/comment-emojis", Handler: rid(middleware.Chain(articleAdminH.EmojiList, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/comment-emojis", Handler: rid(middleware.Chain(articleAdminH.EmojiCreate, gwUser, adminRoles))},
		{Method: http.MethodPut, Path: "/api/v1/admin/comment-emojis/:id", Handler: rid(middleware.Chain(articleAdminH.EmojiUpdate, gwUser, adminRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/comment-emojis/:id", Handler: rid(middleware.Chain(articleAdminH.EmojiDelete, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/article-uploads", Handler: rid(middleware.Chain(articleAdminH.Upload, gwUser, adminRoles))},

		{Method: http.MethodGet, Path: "/api/v1/admin/banners", Handler: rid(middleware.Chain(articleAdminH.ListBanners, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/banners", Handler: rid(middleware.Chain(articleAdminH.CreateBanner, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/banners/upload", Handler: rid(middleware.Chain(articleAdminH.UploadBanner, gwUser, adminRoles))},
		{Method: http.MethodGet, Path: "/api/v1/admin/banners/:id", Handler: rid(middleware.Chain(articleAdminH.GetBanner, gwUser, adminRoles))},
		{Method: http.MethodPut, Path: "/api/v1/admin/banners/:id", Handler: rid(middleware.Chain(articleAdminH.UpdateBanner, gwUser, adminRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/banners/:id", Handler: rid(middleware.Chain(articleAdminH.DeleteBanner, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/shop-uploads", Handler: rid(middleware.Chain(shopUploadH.Upload, gwUser, adminRoles))},

		// 社区文章 — 商家
		{Method: http.MethodGet, Path: "/api/v1/merchant/articles", Handler: rid(middleware.Chain(articleMerchantH.List, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/articles", Handler: rid(middleware.Chain(articleMerchantH.Create, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/articles/:id", Handler: rid(middleware.Chain(articleMerchantH.Detail, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/articles/:id", Handler: rid(middleware.Chain(articleMerchantH.Update, gwShop, merchantRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/articles/:id", Handler: rid(middleware.Chain(articleMerchantH.Delete, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/article-categories", Handler: rid(middleware.Chain(articleMerchantH.CategoryList, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/article-comments", Handler: rid(middleware.Chain(articleMerchantH.CommentList, gwShop, merchantRoles))},
		{Method: http.MethodPatch, Path: "/api/v1/merchant/article-comments/:id", Handler: rid(middleware.Chain(articleMerchantH.CommentPatch, gwShop, merchantRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/article-comments/:id", Handler: rid(middleware.Chain(articleMerchantH.CommentDelete, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/article-uploads", Handler: rid(middleware.Chain(articleMerchantH.Upload, gwShop, merchantRoles))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/notifications/unread-count", Handler: rid(middleware.Chain(notifH.UnreadCount, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/notifications/read-all", Handler: rid(middleware.Chain(notifH.MarkAllRead, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/notifications", Handler: rid(middleware.Chain(notifH.List, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/notifications/:id/read", Handler: rid(middleware.Chain(notifH.MarkRead, gwShop, merchantRoles))},

		// 上传文件：/uploads/products/{shop}/{file}
		{Method: http.MethodGet, Path: "/uploads/products/:shop/:file", Handler: rid(func(w http.ResponseWriter, r *http.Request) {
			p := uploadpath.Abs("products", httpserver.PathParam(r, "shop"), httpserver.PathParam(r, "file"))
			http.ServeFile(w, r, p)
		})},
		{Method: http.MethodGet, Path: "/uploads/exports/:shop/:file", Handler: rid(func(w http.ResponseWriter, r *http.Request) {
			p := uploadpath.Abs("exports", httpserver.PathParam(r, "shop"), httpserver.PathParam(r, "file"))
			http.ServeFile(w, r, p)
		})},
		{Method: http.MethodGet, Path: "/uploads/articles/:shop/:file", Handler: rid(func(w http.ResponseWriter, r *http.Request) {
			p := uploadpath.Abs("articles", httpserver.PathParam(r, "shop"), httpserver.PathParam(r, "file"))
			http.ServeFile(w, r, p)
		})},
		{Method: http.MethodGet, Path: "/uploads/shops/:owner/:file", Handler: rid(func(w http.ResponseWriter, r *http.Request) {
			p := uploadpath.Abs("shops", httpserver.PathParam(r, "owner"), httpserver.PathParam(r, "file"))
			http.ServeFile(w, r, p)
		})},
		{Method: http.MethodGet, Path: "/uploads/reviews/:user/:file", Handler: rid(func(w http.ResponseWriter, r *http.Request) {
			p := uploadpath.Abs("reviews", httpserver.PathParam(r, "user"), httpserver.PathParam(r, "file"))
			http.ServeFile(w, r, p)
		})},
		{Method: http.MethodGet, Path: "/uploads/points-mall/:file", Handler: rid(func(w http.ResponseWriter, r *http.Request) {
			p := uploadpath.Abs("points-mall", httpserver.PathParam(r, "file"))
			http.ServeFile(w, r, p)
		})},
		{Method: http.MethodGet, Path: "/uploads/banners/:file", Handler: rid(func(w http.ResponseWriter, r *http.Request) {
			p := uploadpath.Abs("banners", httpserver.PathParam(r, "file"))
			http.ServeFile(w, r, p)
		})},
	})

	_ = os.MkdirAll(uploadpath.Root(), 0o755)

	go func() {
		logger.Info(fmt.Sprintf("catalog-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
