package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	applog "mymall/pkg/log"
	"mymall/pkg/telemetry"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/data"
	"mymall/services/user-service/internal/handler"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/server"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/uploadpath"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/user-service.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}

	logger, err := applog.New("user-service")
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
		&model.User{},
		&model.SysMenu{},
		&model.SysRole{},
		&model.SysRoleMenu{},
		&model.SysUserRole{},
		&model.SysConfig{},
		&model.UserWallet{},
		&model.UserWalletLog{},
		&model.UserAddress{},
		&model.Region{},
		&model.UserNotification{},
		&model.UserNotificationBatch{},
		&model.UserPoints{},
		&model.UserPointLog{},
		&model.TaskDefinition{},
		&model.UserTaskProgress{},
		&model.UserTaskDedupe{},
	); err != nil {
		log.Fatalf("AutoMigrate 失败：%v", err)
	}

	svcCtx := svc.NewServiceContext(cfg, db)
	if err := svcCtx.Tasks.SeedIfEmpty(); err != nil {
		logger.Warn(fmt.Sprintf("seed tasks failed: %v", err))
	}
	if n, err := svcCtx.Repo.CountRegions(); err == nil && n == 0 {
		if err := svcCtx.Repo.SeedRegionsFromPCA(data.PCACodeJSON); err != nil {
			logger.Warn(fmt.Sprintf("seed regions failed: %v", err))
		} else {
			logger.Info("regions seeded from pca-code.json")
		}
	}
	userLogic := biz.NewUserLogic(context.Background(), svcCtx)
	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(ctx)
	})

	rpcServer := server.StartZRPC(cfg.Server.GRPCPort, userLogic, logger)
	go func() {
		logger.Info(fmt.Sprintf("user-service zRPC 启动 :%d", cfg.Server.GRPCPort))
		rpcServer.Start()
	}()
	defer rpcServer.Stop()

	serverHTTP := httpserver.NewRest(cfg.Server.HTTPPort, cfg.Server.Mode)
	defer serverHTTP.Stop()

	svcCtx.Health = healthReg
	handler.RegisterHandlers(serverHTTP, svcCtx)

	_ = os.MkdirAll(uploadpath.Root(), 0o755)

	go func() {
		logger.Info(fmt.Sprintf("user-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
