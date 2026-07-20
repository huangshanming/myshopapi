package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func AdminResetOwnerPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminResetOwnerPasswordLogic(r.Context(), svcCtx)
		l.AdminResetOwnerPassword(w, r)
	}
}
