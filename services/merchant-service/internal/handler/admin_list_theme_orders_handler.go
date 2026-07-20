package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func AdminListThemeOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminListThemeOrdersLogic(r.Context(), svcCtx)
		l.AdminListThemeOrders(w, r)
	}
}
