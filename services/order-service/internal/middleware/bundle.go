package middleware

import (
	"mymall/pkg/jwt"
	pkgmw "mymall/pkg/middleware"

	"github.com/zeromicro/go-zero/rest"
)

type Bundle struct {
	RequestID                 rest.Middleware
	GatewayIdentity           rest.Middleware
	GatewayIdentityShop       rest.Middleware
	RequireMerchantOwner      rest.Middleware
	RequirePlatformAdmin      rest.Middleware
	RequirePlatformOrMerchant rest.Middleware
}

func NewBundle() Bundle {
	return Bundle{
		RequestID:                 rest.Middleware(pkgmw.RequestID()),
		GatewayIdentity:           rest.Middleware(pkgmw.GatewayIdentity(false)),
		GatewayIdentityShop:       rest.Middleware(pkgmw.GatewayIdentity(true)),
		RequireMerchantOwner:      rest.Middleware(pkgmw.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)),
		RequirePlatformAdmin:      rest.Middleware(pkgmw.RequireRoles(jwt.RolePlatformAdmin)),
		RequirePlatformOrMerchant: rest.Middleware(pkgmw.RequireRoles(jwt.RolePlatformAdmin, jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)),
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
func (b Bundle) Authed(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID, b.GatewayIdentity)(rs)
}
func (b Bundle) MerchantOwner(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID, b.GatewayIdentityShop, b.RequireMerchantOwner)(rs)
}
func (b Bundle) PlatformAdmin(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID, b.GatewayIdentity, b.RequirePlatformAdmin)(rs)
}
func (b Bundle) PlatformOrMerchant(rs []rest.Route) []rest.Route {
	return wrap(b.RequestID, b.GatewayIdentity, b.RequirePlatformOrMerchant)(rs)
}
