package main

import (
	"fmt"
	"log"
	"os"

	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/jwt"
	pkgmiddleware "mymall/pkg/middleware"
	"mymall/services/user-service/internal/handler"
	"mymall/services/user-service/internal/repository"
	"mymall/services/user-service/internal/service"

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

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "user-service"})
	})

	v1 := r.Group("/api/v1")
	v1.Use(pkgmiddleware.ExtractPageReq())
	{
		user := v1.Group("/user")
		user.POST("/login", userHandler.Login)
		user.POST("/register", userHandler.Register)
		user.GET("/profile", jwt.AuthMiddleware(jwtCfg.Secret), userHandler.Profile)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("user-service 启动，监听 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败：%v", err)
	}
}
