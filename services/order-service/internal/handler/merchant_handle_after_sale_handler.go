package handler

import (
	"net/http"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

func MerchantHandleAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMerchantHandleAfterSaleLogic(r.Context(), svcCtx)
		l.MerchantHandleAfterSale(w, r)
	}
}
