package handler

import (
	"net/http"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

func MerchantListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMerchantListLogic(r.Context(), svcCtx)
		l.MerchantList(w, r)
	}
}
