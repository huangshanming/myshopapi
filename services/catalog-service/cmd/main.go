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
	applog "mymall/pkg/log"
	"mymall/pkg/mq"
	"mymall/pkg/telemetry"
	"mymall/pkg/xerr"
	contentlogic "mymall/services/catalog-service/internal/content/logic"
	contentmodel "mymall/services/catalog-service/internal/content/model"
	"mymall/services/catalog-service/internal/handler"
	svcMW "mymall/services/catalog-service/internal/middleware"
	catalogmq "mymall/services/catalog-service/internal/mq"
	notifymodel "mymall/services/catalog-service/internal/notify/model"
	productlogic "mymall/services/catalog-service/internal/product/logic"
	productmodel "mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/server"
	shopopsmodel "mymall/services/catalog-service/internal/shopops/model"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/uploadpath"

	"github.com/redis/go-redis/v9"
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

	handler.RegisterHandlers(serverHTTP, svcCtx, healthReg, svcMW.NewBundle())

	_ = os.MkdirAll(uploadpath.Root(), 0o755)

	go func() {
		logger.Info(fmt.Sprintf("catalog-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
