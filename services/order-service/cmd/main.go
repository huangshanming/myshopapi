package main

// @title           mymall 订单服务 API
// @version         1.0
// @description     下单、查单、取消订单（Saga 异步库存）
// @host            localhost:9080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                JWT Token，格式: Bearer {token}

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	"mymall/pkg/telemetry"
	hadmin "mymall/services/order-service/internal/handler/admin"
	huser "mymall/services/order-service/internal/handler/user"
	hmerchant "mymall/services/order-service/internal/handler/merchant"
	"mymall/services/order-service/internal/model"
	ordermq "mymall/services/order-service/internal/mq"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/order-service.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}

	logger, err := applog.New("order-service")
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
		&model.Order{},
		&model.OrderItem{},
		&model.OrderAfterSale{},
		&model.LogisticsCompany{},
		&model.ProductReview{},
		&model.ProductReviewImage{},
	); err != nil {
		log.Fatalf("AutoMigrate 失败：%v", err)
	}

	svcCtx, err := svc.NewServiceContext(cfg, db)
	if err != nil {
		log.Fatalf("初始化服务依赖失败：%v", err)
	}
	defer svcCtx.Close()

	if svcCtx.UserRPC == nil {
		logger.Warn("user gRPC unavailable")
	}
	if svcCtx.Redis == nil {
		logger.Warn("redis unavailable, order create will reject stock deduct")
	}
	if svcCtx.MQClient == nil {
		logger.Warn("rabbitmq unavailable")
	} else {
		consumer := ordermq.NewConsumer(svcCtx.MQClient, svcCtx.Repo, svcCtx.Redis, svcCtx.UserHTTP, svcCtx.MerchantHTTP, logger)
		if err := consumer.Start(); err != nil {
			logger.Warn("mq consumer start failed")
		}
	}

	orderUser := huser.NewOrderHandler(svcCtx)
	orderMerchant := hmerchant.NewOrderHandler(svcCtx)
	orderAdmin := hadmin.NewOrderHandler(svcCtx)
	reviewUser := huser.NewReviewHandler(svcCtx)
	reviewMerchant := hmerchant.NewReviewHandler(svcCtx)
	reviewAdmin := hadmin.NewReviewHandler(svcCtx)
	logisticsAdmin := hadmin.NewLogisticsHandler(svcCtx)

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(ctx)
	})
	if svcCtx.MQClient != nil {
		healthReg.Register("rabbitmq", svcCtx.MQClient.Ping)
	}
	if svcCtx.Redis != nil {
		healthReg.Register("redis", func(ctx context.Context) error {
			return cache.Ping(ctx, svcCtx.Redis)
		})
	}

	server := httpserver.NewRest(cfg.Server.HTTPPort, cfg.Server.Mode)
	defer server.Stop()

	rid := middleware.RequestID()
	gw := middleware.GatewayIdentity(false)
	gwShop := middleware.GatewayIdentity(true)
	merchantRoles := middleware.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)
	plat := middleware.RequireRoles(jwt.RolePlatformAdmin)
	platOrMerchant := middleware.RequireRoles(jwt.RolePlatformAdmin, jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)

	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: rid(httpserver.Healthz("order-service"))},
		{Method: http.MethodGet, Path: "/readyz", Handler: rid(healthReg.ReadyHandler())},
		{Method: http.MethodGet, Path: "/metrics", Handler: rid(metrics.Handler())},

		{Method: http.MethodPost, Path: "/api/v1/orders", Handler: rid(gw(orderUser.Create))},
		{Method: http.MethodPost, Path: "/api/v1/orders/coupon-preview", Handler: rid(gw(orderUser.CouponPreview))},
		{Method: http.MethodGet, Path: "/api/v1/orders/status-counts", Handler: rid(gw(orderUser.StatusCounts))},
		{Method: http.MethodGet, Path: "/api/v1/orders/after-sales", Handler: rid(gw(orderUser.UserAfterSales))},
		{Method: http.MethodGet, Path: "/api/v1/orders", Handler: rid(gw(orderUser.List))},
		{Method: http.MethodGet, Path: "/api/v1/orders/:id", Handler: rid(gw(orderUser.Detail))},
		{Method: http.MethodPut, Path: "/api/v1/orders/:id/cancel", Handler: rid(gw(orderUser.Cancel))},
		{Method: http.MethodPut, Path: "/api/v1/orders/:id/confirm-receive", Handler: rid(gw(orderUser.ConfirmReceive))},
		{Method: http.MethodPost, Path: "/api/v1/orders/:id/after-sales", Handler: rid(gw(orderUser.CreateAfterSale))},
		{Method: http.MethodGet, Path: "/api/v1/orders/:id/review-eligible", Handler: rid(gw(reviewUser.Eligible))},
		{Method: http.MethodPost, Path: "/api/v1/orders/:id/reviews", Handler: rid(gw(reviewUser.Create))},
		{Method: http.MethodGet, Path: "/api/v1/orders/:id/review", Handler: rid(gw(reviewUser.GetByOrder))},
		{Method: http.MethodPost, Path: "/api/v1/user/review-uploads", Handler: rid(gw(reviewUser.Upload))},
		{Method: http.MethodGet, Path: "/api/v1/products/:id/reviews", Handler: rid(reviewUser.ProductList)},

		{Method: http.MethodGet, Path: "/api/v1/logistics/options", Handler: rid(middleware.Chain(logisticsAdmin.Options, gw, platOrMerchant))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/reviews", Handler: rid(middleware.Chain(reviewMerchant.MerchantList, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/reviews/:id/reply", Handler: rid(middleware.Chain(reviewMerchant.MerchantReply, gwShop, merchantRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/reviews/:id", Handler: rid(middleware.Chain(reviewMerchant.MerchantDelete, gwShop, merchantRoles))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/orders", Handler: rid(middleware.Chain(orderMerchant.MerchantList, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/orders/:id", Handler: rid(middleware.Chain(orderMerchant.MerchantDetail, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/orders/:id/ship", Handler: rid(middleware.Chain(orderMerchant.MerchantShip, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/orders/:id/complete", Handler: rid(middleware.Chain(orderMerchant.MerchantComplete, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/orders/:id/remark", Handler: rid(middleware.Chain(orderMerchant.MerchantRemark, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/after-sales", Handler: rid(middleware.Chain(orderMerchant.MerchantAfterSales, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/after-sales/:id/handle", Handler: rid(middleware.Chain(orderMerchant.MerchantHandleAfterSale, gwShop, merchantRoles))},

		{Method: http.MethodGet, Path: "/api/v1/admin/reviews", Handler: rid(middleware.Chain(reviewAdmin.AdminList, gw, plat))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/reviews/:id", Handler: rid(middleware.Chain(reviewAdmin.AdminDelete, gw, plat))},

		{Method: http.MethodGet, Path: "/api/v1/admin/orders", Handler: rid(middleware.Chain(orderAdmin.AdminList, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/orders/:id", Handler: rid(middleware.Chain(orderAdmin.AdminDetail, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/orders/:id/ship", Handler: rid(middleware.Chain(orderAdmin.AdminShip, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/orders/:id/complete", Handler: rid(middleware.Chain(orderAdmin.AdminComplete, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/orders/:id/remark", Handler: rid(middleware.Chain(orderAdmin.AdminRemark, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/after-sales", Handler: rid(middleware.Chain(orderAdmin.AdminAfterSales, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/after-sales/:id/handle", Handler: rid(middleware.Chain(orderAdmin.AdminHandleAfterSale, gw, plat))},

		{Method: http.MethodGet, Path: "/api/v1/admin/logistics", Handler: rid(middleware.Chain(logisticsAdmin.AdminList, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/logistics", Handler: rid(middleware.Chain(logisticsAdmin.Create, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/logistics/:id", Handler: rid(middleware.Chain(logisticsAdmin.Update, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/logistics/:id/status", Handler: rid(middleware.Chain(logisticsAdmin.UpdateStatus, gw, plat))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/logistics/:id", Handler: rid(middleware.Chain(logisticsAdmin.Delete, gw, plat))},
	})

	go func() {
		logger.Info(fmt.Sprintf("order-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		server.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.Stop()
}
