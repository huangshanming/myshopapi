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

	"mymall/pkg/apidoc"
	"mymall/pkg/cache"
	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/health"
	"mymall/pkg/jwt"
	applog "mymall/pkg/log"
	"mymall/pkg/metrics"
	"mymall/pkg/middleware"
	"mymall/pkg/mq"
	"mymall/pkg/telemetry"
	cataloggrpc "mymall/services/catalog-service/internal/grpc"
	"mymall/services/catalog-service/internal/handler"
	catalogmq "mymall/services/catalog-service/internal/mq"
	"mymall/services/catalog-service/internal/repository"
	"mymall/services/catalog-service/internal/service"

	_ "mymall/services/catalog-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
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

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := service.NewCatalogService(productRepo, categoryRepo, redisClient)
	catalogHandler := handler.NewCatalogHandler(svc)

	if mqClient != nil {
		consumer := catalogmq.NewConsumer(mqClient, svc, logger)
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

	r := gin.Default()
	r.Use(middleware.RequestID())
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "catalog-service"})
	})
	r.GET("/readyz", healthReg.ReadyHandler())
	r.GET("/metrics", metrics.Handler())
	apidoc.MountSwagger(r)

	v1 := r.Group("/api/v1")
	v1.Use(middleware.ExtractPageReq())
	{
		products := v1.Group("/products")
		products.GET("/list", catalogHandler.GetProductList)
		products.GET("/detail", catalogHandler.GetProductDetail)

		categories := v1.Group("/product_category")
		categories.GET("/list", catalogHandler.GetCategoryList)
		categories.GET("/detail", catalogHandler.GetCategoryDetail)

		merchant := v1.Group("/merchant")
		merchant.Use(middleware.GatewayIdentity(true))
		merchant.Use(middleware.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff))
		{
			merchant.GET("/products", catalogHandler.MerchantListProducts)
			merchant.POST("/products", catalogHandler.MerchantCreateProduct)
			merchant.PUT("/products/:id", catalogHandler.MerchantUpdateProduct)
			merchant.PUT("/products/:id/status", catalogHandler.MerchantSetStatus)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.GatewayIdentity(false))
		admin.Use(middleware.RequireRoles(jwt.RolePlatformAdmin))
		{
			admin.GET("/products", catalogHandler.AdminListProducts)
			admin.PUT("/products/:id/off_sale", catalogHandler.AdminForceOffSale)
			admin.POST("/categories", catalogHandler.AdminCreateCategory)
			admin.PUT("/categories/:id", catalogHandler.AdminUpdateCategory)
			admin.DELETE("/categories/:id", catalogHandler.AdminDeleteCategory)
		}
	}

	grpcServer, lis, err := cataloggrpc.Listen(cfg.Server.GRPCPort, svc, logger)
	if err != nil {
		log.Fatalf("gRPC 监听失败：%v", err)
	}
	go func() {
		logger.Info(fmt.Sprintf("catalog-service gRPC 启动 :%d", cfg.Server.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC serve failed")
		}
	}()

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	go func() {
		logger.Info(fmt.Sprintf("catalog-service HTTP 启动 %s", addr))
		if err := r.Run(addr); err != nil {
			log.Fatalf("服务启动失败：%v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	grpcServer.GracefulStop()
}
