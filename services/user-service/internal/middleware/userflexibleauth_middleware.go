package middleware

import (
	"net/http"

	"mymall/pkg/jwt"
	pkgmw "mymall/pkg/middleware"
)

type UserFlexibleAuthMiddleware struct {
	secret string
}

func NewUserFlexibleAuthMiddleware() *UserFlexibleAuthMiddleware {
	// secret filled in ServiceContext after New
	return &UserFlexibleAuthMiddleware{}
}

func NewUserFlexibleAuthMiddlewareWithSecret(secret string) *UserFlexibleAuthMiddleware {
	return &UserFlexibleAuthMiddleware{secret: secret}
}

func (m *UserFlexibleAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	secret := m.secret
	authJWT := jwt.AuthMiddleware(secret)
	authGW := pkgmw.GatewayIdentity(false)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(pkgmw.GatewayUserIDHeader) != "" {
			authGW(next)(w, r)
			return
		}
		authJWT(next)(w, r)
	}
}
