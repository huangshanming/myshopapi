package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func RemarkPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewRemarkPointsOrderLogic(r.Context(), svcCtx)
		l.RemarkPointsOrder(w, r)
	}
}
