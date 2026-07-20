package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func AdminListSlotOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminListSlotOrdersLogic(r.Context(), svcCtx)
		l.AdminListSlotOrders(w, r)
	}
}
