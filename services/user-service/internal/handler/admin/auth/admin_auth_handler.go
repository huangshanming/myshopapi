package auth

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin/auth"
	"mymall/services/user-service/internal/svc"
)

func AuthMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewAuthMeLogic(r.Context(), svcCtx)
		l.AuthMe(w, r)
	}
}
