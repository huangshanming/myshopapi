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

	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	"mymall/pkg/jwt"
	applog "mymall/pkg/log"
	"mymall/pkg/metrics"
	"mymall/pkg/middleware"
	"mymall/pkg/telemetry"
	"mymall/services/order-service/internal/handler"
	"mymall/services/order-service/internal/model"
	ordermq "mymall/services/order-service/internal/mq"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func main() {
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
	if svcCtx.MQClient == nil {
		logger.Warn("rabbitmq unavailable")
	} else {
		consumer := ordermq.NewConsumer(svcCtx.MQClient, svcCtx.Repo, logger)
		if err := consumer.Start(); err != nil {
			logger.Warn("mq consumer start failed")
		}
	}

	orderHandler := handler.NewOrderHandler(svcCtx)

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

	server := httpserver.NewRest(cfg.Server.HTTPPort, cfg.Server.Mode)
	defer server.Stop()

	rid := middleware.RequestID()
	gw := middleware.GatewayIdentity(false)
	gwShop := middleware.GatewayIdentity(true)
	merchantRoles := middleware.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)
	plat := middleware.RequireRoles(jwt.RolePlatformAdmin)

	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: rid(httpserver.Healthz("order-service"))},
		{Method: http.MethodGet, Path: "/readyz", Handler: rid(healthReg.ReadyHandler())},
		{Method: http.MethodGet, Path: "/metrics", Handler: rid(metrics.Handler())},

		{Method: http.MethodPost, Path: "/api/v1/orders", Handler: rid(gw(orderHandler.Create))},
		{Method: http.MethodGet, Path: "/api/v1/orders", Handler: rid(gw(orderHandler.List))},
		{Method: http.MethodGet, Path: "/api/v1/orders/:id", Handler: rid(gw(orderHandler.Detail))},
		{Method: http.MethodPut, Path: "/api/v1/orders/:id/cancel", Handler: rid(gw(orderHandler.Cancel))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/orders", Handler: rid(middleware.Chain(orderHandler.MerchantList, gwShop, merchantRoles))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/orders/:id", Handler: rid(middleware.Chain(orderHandler.MerchantDetail, gwShop, merchantRoles))},

		{Method: http.MethodGet, Path: "/api/v1/admin/orders", Handler: rid(middleware.Chain(orderHandler.AdminList, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/orders/:id", Handler: rid(middleware.Chain(orderHandler.AdminDetail, gw, plat))},
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
