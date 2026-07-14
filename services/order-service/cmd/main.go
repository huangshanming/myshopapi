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

	"mymall/pkg/apidoc"
	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/health"
	"mymall/pkg/jwt"
	applog "mymall/pkg/log"
	"mymall/pkg/metrics"
	"mymall/pkg/middleware"
	"mymall/pkg/mq"
	"mymall/pkg/telemetry"
	grpcclient "mymall/services/order-service/internal/grpc/client"
	"mymall/services/order-service/internal/handler"
	ordermq "mymall/services/order-service/internal/mq"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/service"

	_ "mymall/services/order-service/docs"

	"github.com/gin-gonic/gin"
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

	catalogClient, err := grpcclient.NewCatalogClient(cfg.GRPC.CatalogService)
	if err != nil {
		log.Fatalf("连接 catalog gRPC 失败：%v", err)
	}
	defer catalogClient.Close()

	userClient, err := grpcclient.NewUserClient(cfg.GRPC.UserService)
	if err != nil {
		logger.Warn("user gRPC unavailable")
		userClient = nil
	} else {
		defer userClient.Close()
	}

	var mqClient *mq.Client
	mqc, err := mq.New(cfg.RabbitMQ)
	if err != nil {
		logger.Warn("rabbitmq unavailable")
	} else {
		mqClient = mqc
		defer mqc.Close()
	}

	repo := repository.NewOrderRepository(db)
	var publisher *ordermq.Publisher
	if mqClient != nil {
		publisher = ordermq.NewPublisher(mqClient)
		consumer := ordermq.NewConsumer(mqClient, repo, logger)
		if err := consumer.Start(); err != nil {
			logger.Warn("mq consumer start failed")
		}
	}

	svc := service.NewOrderService(repo, catalogClient, userClient, publisher)
	orderHandler := handler.NewOrderHandler(svc)

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(ctx)
	})
	if mqClient != nil {
		healthReg.Register("rabbitmq", mqClient.Ping)
	}

	r := gin.Default()
	r.Use(middleware.RequestID())
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "order-service"})
	})
	r.GET("/readyz", healthReg.ReadyHandler())
	r.GET("/metrics", metrics.Handler())
	apidoc.MountSwagger(r)

	v1 := r.Group("/api/v1")
	orders := v1.Group("/orders")
	orders.Use(middleware.GatewayIdentity(false))
	{
		orders.POST("", orderHandler.Create)
		orders.GET("", orderHandler.List)
		orders.GET("/:id", orderHandler.Detail)
		orders.PUT("/:id/cancel", orderHandler.Cancel)
	}

	merchant := v1.Group("/merchant")
	merchant.Use(middleware.GatewayIdentity(true))
	merchant.Use(middleware.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff))
	{
		merchant.GET("/orders", orderHandler.MerchantList)
		merchant.GET("/orders/:id", orderHandler.MerchantDetail)
	}

	admin := v1.Group("/admin")
	admin.Use(middleware.GatewayIdentity(false))
	admin.Use(middleware.RequireRoles(jwt.RolePlatformAdmin))
	{
		admin.GET("/orders", orderHandler.AdminList)
		admin.GET("/orders/:id", orderHandler.AdminDetail)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	go func() {
		logger.Info(fmt.Sprintf("order-service HTTP 启动 %s", addr))
		if err := r.Run(addr); err != nil {
			log.Fatalf("服务启动失败：%v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
