package shared

import (
	"net/http"

	"mymall/services/order-service/internal/logic/shared"
	"mymall/services/order-service/internal/svc"
)

func LogisticsOptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shared.NewLogisticsOptionsLogic(r.Context(), svcCtx)
		l.LogisticsOptions(w, r)
	}
}
