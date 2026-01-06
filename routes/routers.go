package routes

import (
	"github.com/gin-gonic/gin"
)
import (
	"github.com/gin-contrib/cors"
	"mymall/middleware"
	"time"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.ExtractPageReq())
	// 第一步：先注册跨域中间件（必须在路由注册之前）
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173",
			"http://127.0.0.1:5173"}, // 前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
	}))

	// 第二步：再注册所有路由（顺序不可颠倒）
	registerAPIRoutes(router)
	registerStaticRoutes(router)

	return router
}

func registerAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			RegisterGoodsRouter(v1)
			RegisterProductsRouter(v1)
			RegisterUserRouter(v1)
		}
	}
}

func registerStaticRoutes(router *gin.Engine) {
	router.Static("/static", "./static")
}
