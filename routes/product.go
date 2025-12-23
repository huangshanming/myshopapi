package routes

import (
	"github.com/gin-gonic/gin"
	"mymall/controllers"
)

func RegisterProductsRouter(router *gin.RouterGroup) {
	// 获取商品列表
	controllers := new(controllers.ProductsController)
	api := router.Group("/products")
	{
		api.GET("/list", controllers.GetList)
	}
}
