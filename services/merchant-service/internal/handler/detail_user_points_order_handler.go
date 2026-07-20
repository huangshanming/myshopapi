package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func DetailUserPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewDetailUserPointsOrderLogic(r.Context(), svcCtx)
		l.DetailUserPointsOrder(w, r)
	}
}
