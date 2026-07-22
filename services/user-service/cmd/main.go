package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mymall/pkg/health"
	applog "mymall/pkg/log"
	"mymall/pkg/telemetry"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/config"
	"mymall/services/user-service/internal/data"
	"mymall/services/user-service/internal/handler"
	"mymall/services/user-service/internal/server"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/uploadpath"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/user-service.yaml"
	}

	var c config.Config
	conf.MustLoad(configPath, &c)
	c.OverlayFromEnv()

	logger, err := applog.New("user-service")
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
	svcCtx := svc.NewServiceContext(&c, sqlConn)
	if err := svcCtx.Tasks.SeedIfEmpty(context.Background()); err != nil {
		logger.Warn(fmt.Sprintf("seed tasks failed: %v", err))
	}
	if n, err := svcCtx.Repo.CountRegions(context.Background()); err == nil && n == 0 {
		if err := svcCtx.Repo.SeedRegionsFromPCA(context.Background(), data.PCACodeJSON); err != nil {
			logger.Warn(fmt.Sprintf("seed regions failed: %v", err))
		} else {
			logger.Info("regions seeded from pca-code.json")
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

	grpcPort := c.GRPCPort()
	rpcServer := server.StartZRPC(grpcPort, c.Etcd.Hosts, svcCtx, logger)
	go func() {
		logger.Info(fmt.Sprintf("user-service zRPC 启动 :%d", grpcPort))
		rpcServer.Start()
	}()
	defer rpcServer.Stop()

	serverHTTP := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer serverHTTP.Stop()

	svcCtx.Health = healthReg
	handler.RegisterHandlers(serverHTTP, svcCtx)

	_ = os.MkdirAll(uploadpath.Root(), 0o755)

	go func() {
		logger.Info(fmt.Sprintf("user-service HTTP(go-zero) 启动 :%d", c.Port))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
