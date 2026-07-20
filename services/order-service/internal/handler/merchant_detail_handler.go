package handler

import (
	"net/http"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

func MerchantDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMerchantDetailLogic(r.Context(), svcCtx)
		l.MerchantDetail(w, r)
	}
}
