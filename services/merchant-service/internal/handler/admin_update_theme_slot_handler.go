package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func AdminUpdateThemeSlotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminUpdateThemeSlotLogic(r.Context(), svcCtx)
		l.AdminUpdateThemeSlot(w, r)
	}
}
