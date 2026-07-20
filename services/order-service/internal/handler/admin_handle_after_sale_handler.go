package handler

import (
	"net/http"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

func AdminHandleAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminHandleAfterSaleLogic(r.Context(), svcCtx)
		l.AdminHandleAfterSale(w, r)
	}
}
