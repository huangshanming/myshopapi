package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mymall/pkg/health"
	applog "mymall/pkg/log"
	"mymall/pkg/xerr"
	biz "mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/config"
	"mymall/services/merchant-service/internal/handler"
	"mymall/services/merchant-service/internal/server"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/merchant-service.yaml"
	}

	var c config.Config
	conf.MustLoad(configPath, &c)
	c.OverlayFromEnv()

	logger, err := applog.New("merchant-service")
	if err != nil {
		log.Fatalf("初始化日志失败：%v", err)
	}
	defer logger.Sync()

	sqlConn := sqlx.NewMysql(c.MySQL.DSN())
	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		rawDB, err := sqlConn.RawDB()
		if err != nil {
			return err
		}
		return rawDB.PingContext(ctx)
	})
	svcCtx := svc.NewServiceContext(&c, sqlConn, healthReg)
	seckillLogic := biz.NewMerchantLogic(svcCtx)
	_, _, _ = seckillLogic.EnsureActiveSession()

	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for range t.C {
			seckillLogic.RotateSeckillSessions()
		}
	}()

	grpcPort := c.GRPCPort()
	rpcServer := server.StartZRPC(grpcPort, svcCtx, logger)
	go func() {
		logger.Info(fmt.Sprintf("merchant-service zRPC 启动 :%d", grpcPort))
		rpcServer.Start()
	}()
	defer rpcServer.Stop()

	serverHTTP := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer serverHTTP.Stop()

	handler.RegisterHandlers(serverHTTP, svcCtx)

	go func() {
		logger.Info(fmt.Sprintf("merchant-service HTTP(go-zero) 启动 :%d", c.Port))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	serverHTTP.Stop()
}
