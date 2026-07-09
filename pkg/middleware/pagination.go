package middleware

import (
	"mymall/pkg/pagination"

	"github.com/gin-gonic/gin"
)

func ExtractPageReq() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pageReq pagination.PageReq
		if err := c.ShouldBindQuery(&pageReq); err != nil {
			_ = c.ShouldBindJSON(&pageReq)
		}
		c.Set("pageReq", &pageReq)
		c.Next()
	}
}
