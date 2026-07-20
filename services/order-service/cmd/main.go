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
	"os"
	"os/signal"
	"syscall"

	"mymall/pkg/cache"
	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	applog "mymall/pkg/log"
	"mymall/pkg/telemetry"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/handler"
	"mymall/services/order-service/internal/model"
	ordermq "mymall/services/order-service/internal/mq"
	"mymall/services/order-service/internal/svc"
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

	svcCtx.Health = healthReg
	handler.RegisterHandlers(server, svcCtx)

	go func() {
		logger.Info(fmt.Sprintf("order-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		server.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.Stop()
}
