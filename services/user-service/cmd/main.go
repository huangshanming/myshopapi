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
	applog "mymall/pkg/log"
	"mymall/pkg/metrics"
	"mymall/pkg/middleware"
	"mymall/pkg/telemetry"
	usergrpc "mymall/services/user-service/internal/grpc"
	"mymall/services/user-service/internal/handler"
	"mymall/services/user-service/internal/repository"
	"mymall/services/user-service/internal/service"

	"mymall/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
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

	jwtCfg := jwt.Config{
		Secret:      cfg.JWT.Secret,
		ConsumerKey: cfg.JWT.ConsumerKey,
		ExpireHours: cfg.JWT.ExpireHours,
		Issuer:      cfg.JWT.Issuer,
	}

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, jwtCfg)
	userHandler := handler.NewUserHandler(svc)

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(ctx)
	})

	r := gin.Default()
	r.Use(middleware.RequestID())
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "user-service"})
	})
	r.GET("/readyz", healthReg.ReadyHandler())
	r.GET("/metrics", metrics.Handler())

	v1 := r.Group("/api/v1")
	v1.Use(middleware.ExtractPageReq())
	{
		user := v1.Group("/user")
		user.POST("/login", userHandler.Login)
		user.POST("/register", userHandler.Register)
		user.GET("/profile", jwt.AuthMiddleware(jwtCfg.Secret), userHandler.Profile)
	}

	grpcServer, lis, err := usergrpc.Listen(cfg.Server.GRPCPort, svc, logger)
	if err != nil {
		log.Fatalf("gRPC 监听失败：%v", err)
	}
	go func() {
		logger.Info(fmt.Sprintf("user-service gRPC 启动 :%d", cfg.Server.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC serve failed")
		}
	}()

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	go func() {
		logger.Info(fmt.Sprintf("user-service HTTP 启动 %s", addr))
		if err := r.Run(addr); err != nil {
			log.Fatalf("服务启动失败：%v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	grpcServer.GracefulStop()
}
