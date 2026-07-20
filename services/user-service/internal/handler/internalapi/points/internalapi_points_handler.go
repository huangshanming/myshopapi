package points

import (
	"net/http"

	"mymall/services/user-service/internal/logic/internalapi/points"
	"mymall/services/user-service/internal/svc"
)

func InternalDeductPointsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points.NewInternalDeductPointsLogic(r.Context(), svcCtx)
		l.InternalDeductPoints(w, r)
	}
}

func InternalRefundPointsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points.NewInternalRefundPointsLogic(r.Context(), svcCtx)
		l.InternalRefundPoints(w, r)
	}
}
