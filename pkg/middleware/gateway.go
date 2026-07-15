package middleware

import (
	"context"
	"net/http"
	"strconv"

	"mymall/pkg/jwt"
	"mymall/pkg/response"
)

type ctxKey string

const (
	GatewayUserIDHeader   = "X-User-Id"
	GatewayUserRoleHeader = "X-User-Role"
	GatewayShopIDHeader   = "X-Shop-Id"

	ctxUserID   ctxKey = "user_id"
	ctxUserRole ctxKey = "user_role"
	ctxShopID   ctxKey = "shop_id"
)

// Middleware go-zero / net/http 中间件签名
type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(h http.HandlerFunc, mws ...Middleware) http.HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func GatewayUserID() Middleware {
	return GatewayIdentity(false)
}

func GatewayIdentity(requireShop bool) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			raw := r.Header.Get(GatewayUserIDHeader)
			if raw == "" {
				response.AbortJSON(w, http.StatusUnauthorized, 401, "missing user id")
				return
			}
			userID, err := strconv.ParseUint(raw, 10, 64)
			if err != nil || userID == 0 {
				response.AbortJSON(w, http.StatusUnauthorized, 401, "invalid user id")
				return
			}
			role := r.Header.Get(GatewayUserRoleHeader)
			if role == "" {
				role = jwt.RoleUser
			}
			if role == jwt.RoleAdmin {
				role = jwt.RolePlatformAdmin
			}
			var shopID uint64
			if shopRaw := r.Header.Get(GatewayShopIDHeader); shopRaw != "" {
				shopID, _ = strconv.ParseUint(shopRaw, 10, 64)
			}
			if requireShop && jwt.IsMerchant(role) && shopID == 0 {
				response.AbortJSON(w, http.StatusForbidden, 403, "missing shop id")
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxUserID, userID)
			ctx = context.WithValue(ctx, ctxUserRole, role)
			ctx = context.WithValue(ctx, ctxShopID, shopID)
			next(w, r.WithContext(ctx))
		}
	}
}

func RequireRoles(roles ...string) Middleware {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
		if r == jwt.RolePlatformAdmin {
			allowed[jwt.RoleAdmin] = struct{}{}
		}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			role := GetUserRole(r.Context())
			if _, ok := allowed[role]; !ok {
				response.AbortJSON(w, http.StatusForbidden, 403, "forbidden")
				return
			}
			next(w, r)
		}
	}
}

func GetUserID(ctx context.Context) (uint64, bool) {
	v := ctx.Value(ctxUserID)
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}

func GetUserRole(ctx context.Context) string {
	v := ctx.Value(ctxUserRole)
	if v == nil {
		return ""
	}
	role, _ := v.(string)
	return role
}

func GetShopID(ctx context.Context) uint64 {
	v := ctx.Value(ctxShopID)
	if v == nil {
		return 0
	}
	id, _ := v.(uint64)
	return id
}
