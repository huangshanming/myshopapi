package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func InternalUnlockCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewInternalUnlockCouponLogic(r.Context(), svcCtx)
		l.InternalUnlockCoupon(w, r)
	}
}
