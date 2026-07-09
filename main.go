package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mymall/db"
	"mymall/middleware"
	"mymall/pkg/config"
	"mymall/routes"
)

func main() {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}

	db.InitGormMySQL(cfg.MySQL)
	defer db.CloseGormMySQL()
	middleware.InitJWT(cfg.JWT)

	router := routes.InitRouter()

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
		if err := router.Run(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gin服务启动失败：%v", err)
		}
	}()
	log.Printf("Gin服务启动成功，监听端口 %d", cfg.Server.HTTPPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")
}
