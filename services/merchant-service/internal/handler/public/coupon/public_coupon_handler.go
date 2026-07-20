package coupon

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/public/coupon"
	"mymall/services/merchant-service/internal/svc"
)

func PublicCouponCenterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewPublicCouponCenterLogic(r.Context(), svcCtx)
		l.PublicCouponCenter(w, r)
	}
}

func PublicCouponPopupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewPublicCouponPopupLogic(r.Context(), svcCtx)
		l.PublicCouponPopup(w, r)
	}
}
