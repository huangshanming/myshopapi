package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const GatewayUserIDHeader = "X-User-Id"

func GatewayUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader(GatewayUserIDHeader)
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "missing user id"})
			return
		}
		userID, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid user id"})
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint64, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}
