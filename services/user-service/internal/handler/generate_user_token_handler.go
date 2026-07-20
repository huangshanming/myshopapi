package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func GenerateUserTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGenerateUserTokenLogic(r.Context(), svcCtx)
		l.GenerateUserToken(w, r)
	}
}
