package middleware

import (
	"net/http"

	"mymall/pkg/jwt"
	pkgmw "mymall/pkg/middleware"
)

type RequirePlatformAdminMiddleware struct{}

func NewRequirePlatformAdminMiddleware() *RequirePlatformAdminMiddleware {
	return &RequirePlatformAdminMiddleware{}
}

func (m *RequirePlatformAdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return pkgmw.RequireRoles(jwt.RolePlatformAdmin)(next)
}
