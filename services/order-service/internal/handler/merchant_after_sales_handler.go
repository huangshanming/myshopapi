package handler

import (
	"net/http"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

func MerchantAfterSalesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMerchantAfterSalesLogic(r.Context(), svcCtx)
		l.MerchantAfterSales(w, r)
	}
}
