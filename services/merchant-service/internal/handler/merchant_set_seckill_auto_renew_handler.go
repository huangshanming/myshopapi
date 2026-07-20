package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantSetSeckillAutoRenewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMerchantSetSeckillAutoRenewLogic(r.Context(), svcCtx)
		l.MerchantSetSeckillAutoRenew(w, r)
	}
}
