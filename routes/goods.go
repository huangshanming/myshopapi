package routes

import (
	"github.com/gin-gonic/gin"
	"mymall/controllers"
	"mymall/middleware"
)

func RegisterGoodsRouter(router *gin.RouterGroup) {
	// 获取商品列表
	controllers := new(controllers.GoodsController)
	api := router.Group("/goods").Use(middleware.JWTAuth())
	{
		api.GET("/list", controllers.GetGoodsList)
	}
}
