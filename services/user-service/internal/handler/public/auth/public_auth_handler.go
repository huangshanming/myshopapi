package auth

import (
	"net/http"

	"mymall/services/user-service/internal/logic/public/auth"
	"mymall/services/user-service/internal/svc"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewLoginLogic(r.Context(), svcCtx)
		l.Login(w, r)
	}
}

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewRegisterLogic(r.Context(), svcCtx)
		l.Register(w, r)
	}
}
