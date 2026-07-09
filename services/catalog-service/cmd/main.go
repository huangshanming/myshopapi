package main

import (
	"fmt"
	"log"
	"os"

	"mymall/pkg/config"
	"mymall/pkg/database"
	pkgmiddleware "mymall/pkg/middleware"
	"mymall/services/catalog-service/internal/handler"
	"mymall/services/catalog-service/internal/repository"
	"mymall/services/catalog-service/internal/service"

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

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := service.NewCatalogService(productRepo, categoryRepo)
	catalogHandler := handler.NewCatalogHandler(svc)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "catalog-service"})
	})

	v1 := r.Group("/api/v1")
	v1.Use(pkgmiddleware.ExtractPageReq())
	{
		products := v1.Group("/products")
		products.GET("/list", catalogHandler.GetProductList)
		products.GET("/detail", catalogHandler.GetProductDetail)

		categories := v1.Group("/product_category")
		categories.GET("/list", catalogHandler.GetCategoryList)
		categories.GET("/detail", catalogHandler.GetCategoryDetail)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("catalog-service 启动，监听 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败：%v", err)
	}
}
