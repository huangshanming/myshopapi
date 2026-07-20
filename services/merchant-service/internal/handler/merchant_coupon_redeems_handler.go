package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantCouponRedeemsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMerchantCouponRedeemsLogic(r.Context(), svcCtx)
		l.MerchantCouponRedeems(w, r)
	}
}
