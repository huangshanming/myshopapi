package middleware

import (
	"net/http"
	"strconv"

	"mymall/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	GatewayUserIDHeader  = "X-User-Id"
	GatewayUserRoleHeader = "X-User-Role"
	GatewayShopIDHeader  = "X-Shop-Id"
)

func GatewayUserID() gin.HandlerFunc {
	return GatewayIdentity(false)
}

// GatewayIdentity 从网关注入的头读取身份；requireShop 要求商家角色必须带 shop_id
func GatewayIdentity(requireShop bool) gin.HandlerFunc {
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

		role := c.GetHeader(GatewayUserRoleHeader)
		if role == "" {
			role = jwt.RoleUser
		}
		if role == jwt.RoleAdmin {
			role = jwt.RolePlatformAdmin
		}
		c.Set("user_role", role)

		var shopID uint64
		if shopRaw := c.GetHeader(GatewayShopIDHeader); shopRaw != "" {
			shopID, _ = strconv.ParseUint(shopRaw, 10, 64)
		}
		c.Set("shop_id", shopID)

		if requireShop && jwt.IsMerchant(role) && shopID == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "msg": "missing shop id"})
			return
		}
		c.Next()
	}
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
		if r == jwt.RolePlatformAdmin {
			allowed[jwt.RoleAdmin] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		roleStr, _ := role.(string)
		if _, ok := allowed[roleStr]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "msg": "forbidden"})
			return
		}
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

func GetUserRole(c *gin.Context) string {
	v, ok := c.Get("user_role")
	if !ok {
		return ""
	}
	role, _ := v.(string)
	return role
}

func GetShopID(c *gin.Context) uint64 {
	v, ok := c.Get("shop_id")
	if !ok {
		return 0
	}
	id, _ := v.(uint64)
	return id
}
