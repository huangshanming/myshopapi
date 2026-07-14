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
	"mymall/pkg/jwt"
	applog "mymall/pkg/log"
	"mymall/pkg/middleware"
	"mymall/services/merchant-service/internal/handler"
	"mymall/services/merchant-service/internal/repository"
	"mymall/services/merchant-service/internal/service"

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

	logger, err := applog.New("merchant-service")
	if err != nil {
		log.Fatalf("初始化日志失败：%v", err)
	}
	defer logger.Sync()

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}

	repo := repository.NewMerchantRepository(db)
	svc := service.NewMerchantService(repo)
	h := handler.NewMerchantHandler(svc)

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
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "merchant-service"})
	})
	r.GET("/readyz", healthReg.ReadyHandler())

	v1 := r.Group("/api/v1")
	{
		merchant := v1.Group("/merchant")
		merchant.Use(middleware.GatewayIdentity(false))
		{
			merchant.POST("/apply", h.Apply)
			merchant.GET("/shops", h.MyShops)
			merchant.PUT("/shops/:id", middleware.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff), h.UpdateMyShop)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.GatewayIdentity(false))
		admin.Use(middleware.RequireRoles(jwt.RolePlatformAdmin))
		{
			admin.GET("/applications", h.AdminListApplications)
			admin.POST("/applications/:id/approve", h.AdminApprove)
			admin.POST("/applications/:id/reject", h.AdminReject)
			admin.GET("/shops", h.AdminListShops)
			admin.GET("/shops/:id", h.AdminGetShop)
			admin.PUT("/shops/:id/disable", h.AdminDisableShop)
			admin.PUT("/shops/:id/enable", h.AdminEnableShop)
		}
	}

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	go func() {
		logger.Info(fmt.Sprintf("merchant-service HTTP 启动 %s", addr))
		if err := r.Run(addr); err != nil {
			log.Fatalf("服务启动失败：%v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
