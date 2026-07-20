package middleware

import (
	"net/http"

	pkgmw "mymall/pkg/middleware"
)

type GatewayIdentityShopMiddleware struct{}

func NewGatewayIdentityShopMiddleware() *GatewayIdentityShopMiddleware {
	return &GatewayIdentityShopMiddleware{}
}

func (m *GatewayIdentityShopMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return pkgmw.GatewayIdentity(true)(next)
}
