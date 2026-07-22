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
	"mymall/pkg/health"
	applog "mymall/pkg/log"
	"mymall/pkg/mq"
	"mymall/pkg/telemetry"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/config"
	contentlogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/handler"
	catalogmq "mymall/services/catalog-service/internal/mq"
	productlogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/server"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/uploadpath"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/catalog-service.yaml"
	}

	var c config.Config
	conf.MustLoad(configPath, &c)
	c.OverlayFromEnv()

	logger, err := applog.New("catalog-service")
	if err != nil {
		log.Fatalf("初始化日志失败：%v", err)
	}
	defer logger.Sync()

	ctx := context.Background()
	shutdownTrace, err := telemetry.Init(ctx, c.Telemetry.ToPkg())
	if err != nil {
		logger.Warn("telemetry init skipped")
	}
	defer shutdownTrace(context.Background())

	sqlConn := sqlx.NewMysql(c.MySQL.DSN())

	var redisClient *redis.Client
	rc, err := cache.NewRedis(c.Redis.ToPkg())
	if err != nil {
		logger.Warn("redis unavailable, cache disabled")
	} else {
		redisClient = rc
	}

	var mqClient *mq.Client
	mqc, err := mq.New(c.RabbitMQ.ToPkg())
	if err != nil {
		logger.Warn("rabbitmq unavailable")
	} else {
		mqClient = mqc
		defer mqc.Close()
	}

	svcCtx := svc.NewServiceContext(&c, sqlConn, redisClient, mqClient)
	if err := svcCtx.ShopRBAC.EnsureShopMenus(context.Background()); err != nil {
		logger.Warn(fmt.Sprintf("seed shop menus: %v", err))
	} else {
		logger.Info("shop menus seeded (layered)")
	}
	catalogLogic := productlogic.NewCatalogLogic(svcCtx)
	productAdminLogic := productlogic.NewProductAdminLogic(svcCtx)
	articleLogic := contentlogic.NewArticleLogic(svcCtx)

	go func() {
		var scheduleMu sync.Mutex
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for range t.C {
			scheduleMu.Lock()
			productAdminLogic.RunSchedules(context.Background())
			articleLogic.RunPublishSchedules(context.Background())
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
		rawDB, err := sqlConn.RawDB()
		if err != nil {
			return err
		}
		return rawDB.PingContext(ctx)
	})
	if redisClient != nil {
		healthReg.Register("redis", func(ctx context.Context) error {
			return cache.Ping(ctx, redisClient)
		})
	}
	if mqClient != nil {
		healthReg.Register("rabbitmq", mqClient.Ping)
	}

	grpcPort := c.GRPCPort()
	rpcServer := server.StartZRPC(grpcPort, catalogLogic, logger)
	go func() {
		logger.Info(fmt.Sprintf("catalog-service zRPC 启动 :%d", grpcPort))
		rpcServer.Start()
	}()
	defer rpcServer.Stop()

	serverHTTP := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer serverHTTP.Stop()

	svcCtx.Health = healthReg
	handler.RegisterHandlers(serverHTTP, svcCtx)

	_ = os.MkdirAll(uploadpath.Root(), 0o755)

	go func() {
		logger.Info(fmt.Sprintf("catalog-service HTTP(go-zero) 启动 :%d", c.Port))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
