1. 项目目录结构
project/
├── main.go
├── go.mod
├── go.sum
├── config/
│   └── config.go
├── routes/
│   ├── routes.go          # 路由注册入口
│   ├── middleware.go      # 中间件
│   ├── user.go           # 用户路由模块
│   ├── goods.go          # 商品路由模块
│   └── order.go          # 订单路由模块
├── controllers/
│   ├── user.go
│   ├── goods.go
│   └── order.go
├── services/
│   ├── user.go
│   ├── goods.go
│   └── order.go
├── models/
│   ├── user.go
│   ├── goods.go
│   └── order.go
└── utils/
    └── response.go


2. 路由注册入口 (routes/routes.go)
package routes

import (
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	router := gin.Default()

	// 全局中间件
	router.Use(LoggerMiddleware(), RecoveryMiddleware(), CorsMiddleware())

	// 注册API路由
	registerAPIRoutes(router)

	// 注册静态文件路由
	registerStaticRoutes(router)

	return router
}

// registerAPIRoutes 注册API路由
func registerAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// v1版本路由
		v1 := api.Group("/v1")
		{
			// 各模块路由注册
			RegisterUserRoutes(v1)
			RegisterGoodsRoutes(v1)
			RegisterOrderRoutes(v1)
		}

		// 未来可以添加v2版本
		// v2 := api.Group("/v2")
		// {
		//     RegisterUserRoutesV2(v2)
		// }
	}
}

// registerStaticRoutes 注册静态文件路由
func registerStaticRoutes(router *gin.Engine) {
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
}

3. 中间件设计 (routes/middleware.go)
package routes

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录日志
		duration := time.Since(start)
		fmt.Printf("[%s] %s %s %v\n",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			duration,
		)
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(500, gin.H{
					"code":    500,
					"message": "服务器内部错误",
					"error":   fmt.Sprintf("%v", err),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CorsMiddleware 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "未授权访问"})
			c.Abort()
			return
		}

		// 验证token逻辑
		// userID, err := utils.ValidateToken(token)
		// if err != nil {
		//     c.JSON(401, gin.H{"error": "Token无效"})
		//     c.Abort()
		//     return
		// }

		// c.Set("userID", userID)
		c.Next()
	}
}

4. 用户路由模块 (routes/user.go)
package routes

import (
	"project/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(router *gin.RouterGroup) {
	userController := controllers.NewUserController()

	user := router.Group("/users")
	{
		// 公开路由（不需要认证）
		user.POST("/register", userController.Register)
		user.POST("/login", userController.Login)
		user.GET("/:id", userController.GetUserInfo)

		// 需要认证的路由
		auth := user.Group("")
		auth.Use(AuthMiddleware()) // 应用认证中间件
		{
			auth.PUT("/:id", userController.UpdateUser)
			auth.DELETE("/:id", userController.DeleteUser)
			auth.GET("/profile", userController.GetProfile)
		}

		// 管理员路由
		admin := user.Group("")
		admin.Use(AuthMiddleware(), AdminMiddleware()) // 认证+管理员权限
		{
			admin.GET("", userController.GetUserList)
			admin.PUT("/:id/status", userController.UpdateUserStatus)
		}
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否有管理员权限
		// userID := c.GetString("userID")
		// if !utils.IsAdmin(userID) {
		//     c.JSON(403, gin.H{"error": "权限不足"})
		//     c.Abort()
		//     return
		// }
		c.Next()
	}
}

5. 商品路由模块 (routes/goods.go)
package routes

import (
	"project/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterGoodsRoutes 注册商品相关路由
func RegisterGoodsRoutes(router *gin.RouterGroup) {
	goodsController := controllers.NewGoodsController()

	goods := router.Group("/goods")
	{
		// 公开路由
		goods.GET("", goodsController.GetGoodsList)
		goods.GET("/:id", goodsController.GetGoodsDetail)
		goods.GET("/search", goodsController.SearchGoods)

		// 需要认证的路由
		auth := goods.Group("")
		auth.Use(AuthMiddleware())
		{
			auth.POST("", goodsController.CreateGoods)
			auth.PUT("/:id", goodsController.UpdateGoods)
			auth.DELETE("/:id", goodsController.DeleteGoods)
			auth.POST("/:id/like", goodsController.LikeGoods)
		}
	}

	// 商品分类路由
	category := router.Group("/categories")
	{
		category.GET("", goodsController.GetCategories)
		category.GET("/:id/goods", goodsController.GetGoodsByCategory)
	}
}

6. 订单路由模块 (routes/order.go)
package routes

import (
	"project/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterOrderRoutes 注册订单相关路由
func RegisterOrderRoutes(router *gin.RouterGroup) {
	orderController := controllers.NewOrderController()

	order := router.Group("/orders")
	order.Use(AuthMiddleware()) // 所有订单路由都需要认证
	{
		order.POST("", orderController.CreateOrder)
		order.GET("", orderController.GetOrderList)
		order.GET("/:id", orderController.GetOrderDetail)
		order.PUT("/:id/cancel", orderController.CancelOrder)
		order.PUT("/:id/status", orderController.UpdateOrderStatus)
		order.GET("/:id/payment", orderController.GetPaymentInfo)
	}
}

7. 控制器示例 (controllers/user.go)
package controllers

import (
	"net/http"

	"project/services"
	"project/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// Register 用户注册
func (ctrl *UserController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	user, err := ctrl.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, "注册成功", user)
}

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	token, user, err := ctrl.userService.Login(req.Username, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	utils.SuccessResponse(c, "登录成功", gin.H{
		"token": token,
		"user":  user,
	})
}

// GetUserInfo 获取用户信息
func (ctrl *UserController) GetUserInfo(c *gin.Context) {
	userID := c.Param("id")
	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	utils.SuccessResponse(c, "获取成功", user)
}

// 其他方法...

8. 响应工具类 (utils/response.go)
package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

9. 主程序入口 (main.go)
package main

import (
	"project/routes"
)

func main() {
	// 初始化路由
	router := routes.InitRouter()

	// 启动服务器
	router.Run(":8080")
}

设计优势
模块化设计：每个业务模块独立，便于维护和扩展

分层架构：路由层、控制器层、服务层、模型层分离

中间件机制：统一的认证、日志、异常处理

版本控制：支持API版本管理

权限控制：灵活的权限中间件

统一响应格式：标准化的API响应结构

这种设计让项目结构清晰，易于维护和扩展，适合中大型项目的开发。