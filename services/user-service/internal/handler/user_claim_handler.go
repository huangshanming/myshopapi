package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func UserClaimHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUserClaimLogic(r.Context(), svcCtx)
		l.UserClaim(w, r)
	}
}
