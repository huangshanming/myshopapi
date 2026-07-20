package middleware

import (
	"net/http"

	"mymall/pkg/jwt"
	pkgmw "mymall/pkg/middleware"
)

type RequireMerchantOwnerMiddleware struct{}

func NewRequireMerchantOwnerMiddleware() *RequireMerchantOwnerMiddleware {
	return &RequireMerchantOwnerMiddleware{}
}

func (m *RequireMerchantOwnerMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return pkgmw.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)(next)
}
