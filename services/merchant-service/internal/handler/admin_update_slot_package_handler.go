package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func AdminUpdateSlotPackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminUpdateSlotPackageLogic(r.Context(), svcCtx)
		l.AdminUpdateSlotPackage(w, r)
	}
}
