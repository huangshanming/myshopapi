package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func ResetAdminPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewResetAdminPasswordLogic(r.Context(), svcCtx)
		l.ResetAdminPassword(w, r)
	}
}
