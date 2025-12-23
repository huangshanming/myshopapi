package main

import "fmt"
import "mymall/routes"
import "mymall/db"
import "os"
import "log"
import "net/http"
import "os/signal"
import "syscall"

func main() {
	db.InitGormMySQL()
	defer db.CloseGormMySQL()
	router := routes.InitRouter()

	fmt.Println("Hello World")

	// 3. 启动服务与优雅关机
	go func() {
		if err := router.Run(":8087"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gin服务启动失败：%v", err)
		}
	}()
	log.Println("Gin服务启动成功，监听端口8080")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

}
