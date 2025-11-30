package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	// 注册API路由
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
		}
	}
}

func registerStaticRoutes(router *gin.Engine) {
	router.Static("/static", "./static")
}
