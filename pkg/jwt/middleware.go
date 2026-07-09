package jwt

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "请求头缺少Authorization字段",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Authorization格式错误（正确格式：Bearer <token>）",
			})
			c.Abort()
			return
		}

		claims, err := ParseToken(parts[1], secret)
		if err != nil {
			msg := "Token校验失败"
			if errors.Is(err, ErrExpiredToken) {
				msg = "Token已过期"
			} else if errors.Is(err, ErrInvalidToken) {
				msg = "Token无效"
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  msg,
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)
		c.Next()
	}
}
