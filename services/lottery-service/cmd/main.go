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
	"mymall/pkg/health"
	applog "mymall/pkg/log"
	"mymall/pkg/xerr"
	lotcfg "mymall/services/lottery-service/internal/config"
	"mymall/services/lottery-service/internal/handler"
	"mymall/services/lottery-service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/lottery-service.yaml"
	}

	var c lotcfg.Config
	conf.MustLoad(configPath, &c)
	c.OverlayFromEnv()

	logger, err := applog.New("lottery-service")
	if err != nil {
		log.Fatalf("初始化日志失败：%v", err)
	}
	defer logger.Sync()

	sqlConn := sqlx.NewMysql(c.MySQL.DSN())
	rdb, err := cache.NewRedis(config.RedisConfig{
		Host:     c.Redis.Host,
		Port:     c.Redis.Port,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})
	if err != nil {
		log.Fatalf("连接 Redis 失败：%v", err)
	}
	defer rdb.Close()

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		rawDB, err := sqlConn.RawDB()
		if err != nil {
			return err
		}
		return rawDB.PingContext(ctx)
	})
	healthReg.Register("redis", func(ctx context.Context) error {
		return cache.Ping(ctx, rdb)
	})
	svcCtx := svc.NewServiceContext(&c, sqlConn, rdb, healthReg)

	serverHTTP := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer serverHTTP.Stop()
	handler.RegisterHandlers(serverHTTP, svcCtx)

	go func() {
		logger.Info(fmt.Sprintf("lottery-service HTTP 启动 :%d", c.Port))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
