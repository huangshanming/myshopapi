package routes

import (
	"github.com/gin-gonic/gin"
	"mymall/controllers"
)

func RegisterProductsRouter(router *gin.RouterGroup) {
	// 获取商品列表
	cate_controllers := new(controllers.ProductCategoryController)
	cate_api := router.Group("/product_category")
	{
		cate_api.GET("/list", cate_controllers.GetList)
		cate_api.GET("/detail", cate_controllers.GetDetail)
	}
	product_controllers := new(controllers.ProductsController)
	api := router.Group("/products")
	{
		api.GET("/list", product_controllers.GetList)
		api.GET("/detail", product_controllers.GetDetail)
	}
}
