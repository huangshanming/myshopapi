package routes

func RegisterGoodsRouter(router *gin.Engine) {
	// 获取商品列表
	api := router.GROUP("/goods") {
		api.GET("/list", controllers.GetGoodsList)
	}
}
