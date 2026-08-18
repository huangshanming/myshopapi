package middleware

import (
	"net/http"

	pkgmw "mymall/pkg/middleware"
)

type GatewayIdentityMiddleware struct{}

func NewGatewayIdentityMiddleware() *GatewayIdentityMiddleware { return &GatewayIdentityMiddleware{} }

func (m *GatewayIdentityMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return pkgmw.GatewayIdentity(false)(next)
}
