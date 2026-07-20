package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func ListUserPointsOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewListUserPointsOrdersLogic(r.Context(), svcCtx)
		l.ListUserPointsOrders(w, r)
	}
}
