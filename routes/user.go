package routes

import (
	"github.com/gin-gonic/gin"
	"mymall/controllers"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	// 获取商品列表
	user_controllers := new(controllers.UsersController)
	api := router.Group("/user")
	{
		api.POST("/login", user_controllers.Login)
		api.POST("/register", user_controllers.Register)
	}
}
