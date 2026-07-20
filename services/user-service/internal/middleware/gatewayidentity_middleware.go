package middleware

import (
	pkgmw "mymall/pkg/middleware"
	"net/http"
)

type GatewayIdentityMiddleware struct{}

func NewGatewayIdentityMiddleware() *GatewayIdentityMiddleware { return &GatewayIdentityMiddleware{} }

func (m *GatewayIdentityMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return pkgmw.GatewayIdentity(false)(next)
}
