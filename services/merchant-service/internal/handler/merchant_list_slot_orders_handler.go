package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantListSlotOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMerchantListSlotOrdersLogic(r.Context(), svcCtx)
		l.MerchantListSlotOrders(w, r)
	}
}
