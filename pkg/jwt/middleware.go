package jwt

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"mymall/pkg/response"
)

type ctxKey string

const ctxClaims ctxKey = "jwt_claims"

// AuthMiddleware 本地 JWT 校验（直连服务时使用；经 APISIX 时用 GatewayIdentity）
func AuthMiddleware(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.AbortJSON(w, http.StatusUnauthorized, 401, "请求头缺少Authorization字段")
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.AbortJSON(w, http.StatusUnauthorized, 401, "Authorization格式错误（正确格式：Bearer <token>）")
				return
			}
			claims, err := ParseToken(parts[1], secret)
			if err != nil {
				msg := "Token校验失败"
				if errors.Is(err, ErrExpiredToken) {
					msg = "Token已过期"
				} else if errors.Is(err, ErrInvalidToken) {
					msg = "Token无效"
				}
				response.AbortJSON(w, http.StatusUnauthorized, 401, msg)
				return
			}
			ctx := context.WithValue(r.Context(), ctxClaims, claims)
			next(w, r.WithContext(ctx))
		}
	}
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	v := ctx.Value(ctxClaims)
	if v == nil {
		return nil, false
	}
	c, ok := v.(*Claims)
	return c, ok
}
