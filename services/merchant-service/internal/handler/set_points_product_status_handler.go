package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func SetPointsProductStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewSetPointsProductStatusLogic(r.Context(), svcCtx)
		l.SetPointsProductStatus(w, r)
	}
}
