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
	"mymall/pkg/mq"
	"mymall/pkg/telemetry"
	"mymall/services/catalog-service/internal/handler"
	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/model"
	catalogmq "mymall/services/catalog-service/internal/mq"
	"mymall/services/catalog-service/internal/server"
	"mymall/services/catalog-service/internal/svc"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
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
		&model.Product{},
		&model.ProductCategory{},
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
	catalogLogic := logic.NewCatalogLogic(svcCtx)
	catalogHandler := handler.NewCatalogHandler(svcCtx)

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
	merchantRoles := middleware.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)
	adminRoles := middleware.RequireRoles(jwt.RolePlatformAdmin)

	serverHTTP.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: rid(httpserver.Healthz("catalog-service"))},
		{Method: http.MethodGet, Path: "/readyz", Handler: rid(healthReg.ReadyHandler())},
		{Method: http.MethodGet, Path: "/metrics", Handler: rid(metrics.Handler())},

		{Method: http.MethodGet, Path: "/api/v1/products/list", Handler: rid(catalogHandler.GetProductList)},
		{Method: http.MethodGet, Path: "/api/v1/products/detail", Handler: rid(catalogHandler.GetProductDetail)},
		{Method: http.MethodGet, Path: "/api/v1/product_category/list", Handler: rid(catalogHandler.GetCategoryList)},
		{Method: http.MethodGet, Path: "/api/v1/product_category/detail", Handler: rid(catalogHandler.GetCategoryDetail)},

		{Method: http.MethodGet, Path: "/api/v1/merchant/products", Handler: rid(middleware.Chain(catalogHandler.MerchantListProducts, gwShop, merchantRoles))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products", Handler: rid(middleware.Chain(catalogHandler.MerchantCreateProduct, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/products/:id", Handler: rid(middleware.Chain(catalogHandler.MerchantUpdateProduct, gwShop, merchantRoles))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/products/:id/status", Handler: rid(middleware.Chain(catalogHandler.MerchantSetStatus, gwShop, merchantRoles))},

		{Method: http.MethodGet, Path: "/api/v1/admin/products", Handler: rid(middleware.Chain(catalogHandler.AdminListProducts, gwUser, adminRoles))},
		{Method: http.MethodPut, Path: "/api/v1/admin/products/:id/off_sale", Handler: rid(middleware.Chain(catalogHandler.AdminForceOffSale, gwUser, adminRoles))},
		{Method: http.MethodPost, Path: "/api/v1/admin/categories", Handler: rid(middleware.Chain(catalogHandler.AdminCreateCategory, gwUser, adminRoles))},
		{Method: http.MethodPut, Path: "/api/v1/admin/categories/:id", Handler: rid(middleware.Chain(catalogHandler.AdminUpdateCategory, gwUser, adminRoles))},
		{Method: http.MethodDelete, Path: "/api/v1/admin/categories/:id", Handler: rid(middleware.Chain(catalogHandler.AdminDeleteCategory, gwUser, adminRoles))},
	})

	go func() {
		logger.Info(fmt.Sprintf("catalog-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
