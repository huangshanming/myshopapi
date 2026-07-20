package middleware

import (
	"mymall/pkg/jwt"
	pkgmw "mymall/pkg/middleware"
	"net/http"
)

type RequirePlatformAdminMiddleware struct{}

func NewRequirePlatformAdminMiddleware() *RequirePlatformAdminMiddleware {
	return &RequirePlatformAdminMiddleware{}
}

func (m *RequirePlatformAdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return pkgmw.RequireRoles(jwt.RolePlatformAdmin)(next)
}
