package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func ShipPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewShipPointsOrderLogic(r.Context(), svcCtx)
		l.ShipPointsOrder(w, r)
	}
}
