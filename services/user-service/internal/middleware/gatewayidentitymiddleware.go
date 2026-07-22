// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package middleware

import "net/http"

type GatewayIdentityMiddleware struct {
}

func NewGatewayIdentityMiddleware() *GatewayIdentityMiddleware {
	return &GatewayIdentityMiddleware{}
}

func (m *GatewayIdentityMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}
