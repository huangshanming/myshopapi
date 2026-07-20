package middleware

import (
	"mymall/pkg/jwt"
	pkgmw "mymall/pkg/middleware"

	"github.com/zeromicro/go-zero/rest"
)

// Bundle holds named middlewares declared in api/*.api @server groups.
type Bundle struct {
	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	RequireMerchantOwner rest.Middleware
	RequirePlatformAdmin rest.Middleware
}

func NewBundle() Bundle {
	return Bundle{
		RequestID:            rest.Middleware(pkgmw.RequestID()),
		GatewayIdentity:      rest.Middleware(pkgmw.GatewayIdentity(false)),
		RequireMerchantOwner: rest.Middleware(pkgmw.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)),
		RequirePlatformAdmin: rest.Middleware(pkgmw.RequireRoles(jwt.RolePlatformAdmin)),
	}
}

func wrap(ms ...rest.Middleware) func([]rest.Route) []rest.Route {
	return func(rs []rest.Route) []rest.Route {
		return rest.WithMiddlewares(ms, rs...)
	}
}

// Public applies RequestID (api: middleware: RequestID).
func (b Bundle) Public(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID)(rs)
}

// Authed applies RequestID + GatewayIdentity.
func (b Bundle) Authed(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID, b.GatewayIdentity)(rs)
}

// MerchantOwner applies RequestID + GatewayIdentity + RequireMerchantOwner.
func (b Bundle) MerchantOwner(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID, b.GatewayIdentity, b.RequireMerchantOwner)(rs)
}

// PlatformAdmin applies RequestID + GatewayIdentity + RequirePlatformAdmin.
func (b Bundle) PlatformAdmin(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID, b.GatewayIdentity, b.RequirePlatformAdmin)(rs)
}
