package logistics

import (
	"net/http"

	"mymall/services/order-service/internal/logic/shared/logistics"
	"mymall/services/order-service/internal/svc"
)

func LogisticsOptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logistics.NewLogisticsOptionsLogic(r.Context(), svcCtx)
		l.LogisticsOptions(w, r)
	}
}
