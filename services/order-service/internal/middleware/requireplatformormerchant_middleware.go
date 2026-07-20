package middleware

import (
	"net/http"

	"mymall/pkg/jwt"
	pkgmw "mymall/pkg/middleware"
)

type RequirePlatformOrMerchantMiddleware struct{}

func NewRequirePlatformOrMerchantMiddleware() *RequirePlatformOrMerchantMiddleware {
	return &RequirePlatformOrMerchantMiddleware{}
}

func (m *RequirePlatformOrMerchantMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return pkgmw.RequireRoles(jwt.RolePlatformAdmin, jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)(next)
}
