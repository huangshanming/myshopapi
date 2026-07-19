package main

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
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	canalclient "mymall/services/inventory-sync-service/internal/canal"
	"mymall/services/inventory-sync-service/internal/handler"
	"mymall/services/inventory-sync-service/internal/preheat"
	"mymall/services/inventory-sync-service/internal/svc"
	"mymall/services/inventory-sync-service/internal/sync"

	"go.uber.org/zap"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/inventory-sync-service.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}

	logger, err := applog.New("inventory-sync-service")
	if err != nil {
		log.Fatalf("初始化日志失败：%v", err)
	}
	defer logger.Sync()

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}

	rdb, err := cache.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("连接 Redis 失败：%v", err)
	}
	defer rdb.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, err := preheat.LoadAllSkuStock(ctx, db, rdb, logger)
	if err != nil {
		log.Fatalf("库存预热失败：%v", err)
	}
	logger.Info("sku stock preheated",
		zap.Int("scanned", st.Scanned),
		zap.Int("filled", st.Filled),
		zap.Int("pulled_down", st.PulledDown),
		zap.Int("kept", st.Kept),
	)

	syncHandler := sync.NewHandler(rdb, logger)
	consumer := canalclient.NewConsumer(cfg.Canal, syncHandler, logger)
	go func() {
		if err := consumer.Run(ctx); err != nil && err != context.Canceled {
			logger.Warn("canal consumer stopped", zap.Error(err))
		}
	}()

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(c context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(c)
	})
	healthReg.Register("redis", func(c context.Context) error {
		return cache.Ping(c, rdb)
	})

	svcCtx := svc.NewServiceContext(db, rdb, healthReg)

	server := httpserver.NewRest(cfg.Server.HTTPPort, cfg.Server.Mode)
	defer server.Stop()
	rid := middleware.RequestID()
	handler.RegisterHandlers(server, svcCtx, rid)

	go func() {
		logger.Info(fmt.Sprintf("inventory-sync-service HTTP(goctl) 启动 :%d", cfg.Server.HTTPPort))
		server.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
	server.Stop()
}
