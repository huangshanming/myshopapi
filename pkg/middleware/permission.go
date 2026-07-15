package middleware

import (
	"net/http"

	"mymall/pkg/jwt"
	"mymall/pkg/response"
)

// PermChecker 由业务层注入（查 DB 权限）
type PermChecker interface {
	IsSuperAdmin(userID uint64) bool
	HasPerm(userID uint64, code string) bool
}

// RequirePermission 要求具备指定权限码；platform_admin 粗粒度不足时以 perms 为准。
// 须在 GatewayIdentity 或等价注入 user_id 之后使用。
func RequirePermission(checker PermChecker, code string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r.Context())
			if !ok || userID == 0 {
				response.AbortJSON(w, http.StatusUnauthorized, 401, "unauthorized")
				return
			}
			role := GetUserRole(r.Context())
			if !jwt.IsPlatformAdmin(role) {
				response.AbortJSON(w, http.StatusForbidden, 403, "forbidden")
				return
			}
			if code == "" || checker == nil {
				next(w, r)
				return
			}
			if checker.IsSuperAdmin(userID) || checker.HasPerm(userID, code) {
				next(w, r)
				return
			}
			response.AbortJSON(w, http.StatusForbidden, 403, "无权限")
		}
	}
}
