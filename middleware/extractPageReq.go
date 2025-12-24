package middleware

import (
	"github.com/gin-gonic/gin"
	"mymall/common"
)

func ExtractPageReq() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pageReq common.PageReq
		// 优先绑定GET查询参数，若失败则绑定JSON请求体（适配多种请求方式）
		if err := c.ShouldBindQuery(&pageReq); err != nil {
			_ = c.ShouldBindJSON(&pageReq)
		}

		// 2. 将分页参数存入Gin上下文，供后续接口使用
		c.Set("pageReq", &pageReq)

		// 3. 继续执行后续中间件和接口逻辑
		c.Next()
	}
}
