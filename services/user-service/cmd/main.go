package main

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
	"mymall/services/user-service/internal/handler"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/server"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func main() {
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
	if err := database.AutoMigrateIfDebug(cfg.Server.Mode, db, &model.User{}); err != nil {
		log.Fatalf("AutoMigrate 失败：%v", err)
	}

	svcCtx := svc.NewServiceContext(cfg, db)
	userLogic := logic.NewUserLogic(svcCtx)
	userHandler := handler.NewUserHandler(svcCtx)

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

	rid := middleware.RequestID()
	authJWT := jwt.AuthMiddleware(svcCtx.JWT.Secret)
	authGW := middleware.GatewayIdentity(false)

	profile := func(w http.ResponseWriter, r *http.Request) {
		// 优先网关头；无头时回退 JWT（本地直连）
		if r.Header.Get(middleware.GatewayUserIDHeader) != "" {
			authGW(userHandler.Profile)(w, r)
			return
		}
		authJWT(userHandler.Profile)(w, r)
	}

	serverHTTP.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: rid(httpserver.Healthz("user-service"))},
		{Method: http.MethodGet, Path: "/readyz", Handler: rid(healthReg.ReadyHandler())},
		{Method: http.MethodGet, Path: "/metrics", Handler: rid(metrics.Handler())},

		{Method: http.MethodPost, Path: "/api/v1/user/login", Handler: rid(userHandler.Login)},
		{Method: http.MethodPost, Path: "/api/v1/user/register", Handler: rid(userHandler.Register)},
		{Method: http.MethodGet, Path: "/api/v1/user/profile", Handler: rid(profile)},
	})

	go func() {
		logger.Info(fmt.Sprintf("user-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
