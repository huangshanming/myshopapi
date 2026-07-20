package middleware

import (
	"net/http"

	"mymall/pkg/jwt"
	pkgmw "mymall/pkg/middleware"

	"github.com/zeromicro/go-zero/rest"
)

type Bundle struct {
	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	RequirePlatformAdmin rest.Middleware
	UserFlexibleAuth     rest.Middleware
	jwtSecret            string
}

func NewBundle(jwtSecret string) Bundle {
	authJWT := jwt.AuthMiddleware(jwtSecret)
	authGW := pkgmw.GatewayIdentity(false)
	flex := rest.Middleware(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(pkgmw.GatewayUserIDHeader) != "" {
				authGW(next)(w, r)
				return
			}
			authJWT(next)(w, r)
		}
	})
	return Bundle{
		RequestID:            rest.Middleware(pkgmw.RequestID()),
		GatewayIdentity:      rest.Middleware(authGW),
		RequirePlatformAdmin: rest.Middleware(pkgmw.RequireRoles(jwt.RolePlatformAdmin)),
		UserFlexibleAuth:     flex,
		jwtSecret:            jwtSecret,
	}
}

func wrap(ms ...rest.Middleware) func([]rest.Route) []rest.Route {
	return func(rs []rest.Route) []rest.Route {
		return rest.WithMiddlewares(ms, rs...)
	}
}

func (b Bundle) Public(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID)(rs)
}
func (b Bundle) UserFlexible(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID, b.UserFlexibleAuth)(rs)
}
func (b Bundle) PlatformAdmin(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID, b.GatewayIdentity, b.RequirePlatformAdmin)(rs)
}
