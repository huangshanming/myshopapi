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
	"mymall/pkg/health"
	applog "mymall/pkg/log"
	"mymall/pkg/telemetry"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/config"
	"mymall/services/order-service/internal/handler"
	ordermq "mymall/services/order-service/internal/mq"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/order-service.yaml"
	}

	var c config.Config
	conf.MustLoad(configPath, &c)
	c.OverlayFromEnv()

	logger, err := applog.New("order-service")
	if err != nil {
		log.Fatalf("初始化日志失败：%v", err)
	}
	defer logger.Sync()

	ctx := context.Background()
	shutdownTrace, err := telemetry.Init(ctx, c.AppTelemetry.ToPkg())
	if err != nil {
		logger.Warn("telemetry init skipped")
	}
	defer shutdownTrace(context.Background())

	sqlConn := sqlx.NewMysql(c.MySQL.DSN())
	svcCtx, err := svc.NewServiceContext(&c, sqlConn)
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
		consumer := ordermq.NewConsumer(svcCtx.MQClient, svcCtx.Repo, svcCtx.Redis, svcCtx.UserRPC, svcCtx.MerchantRPC, logger)
		if err := consumer.Start(); err != nil {
			logger.Warn("mq consumer start failed")
		}
	}

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		rawDB, err := sqlConn.RawDB()
		if err != nil {
			return err
		}
		return rawDB.PingContext(ctx)
	})
	if svcCtx.MQClient != nil {
		healthReg.Register("rabbitmq", svcCtx.MQClient.Ping)
	}
	if svcCtx.Redis != nil {
		healthReg.Register("redis", func(ctx context.Context) error {
			return cache.Ping(ctx, svcCtx.Redis)
		})
	}

	serverHTTP := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer serverHTTP.Stop()

	svcCtx.Health = healthReg
	handler.RegisterHandlers(serverHTTP, svcCtx)

	go func() {
		logger.Info(fmt.Sprintf("order-service HTTP(go-zero) 启动 :%d", c.Port))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	serverHTTP.Stop()
}
