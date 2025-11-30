package routes

import (
	"github.com/gin-gonic/gin"
	"mymall/controllers"
)

func RegisterGoodsRouter(router *gin.RouterGroup) {
	// 获取商品列表
	controllers := new(controllers.GoodsController)
	api := router.Group("/goods")
	{
		api.GET("/list", controllers.GetGoodsList)
	}
}
