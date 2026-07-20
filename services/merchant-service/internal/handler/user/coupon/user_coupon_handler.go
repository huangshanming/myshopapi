package coupon

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/user/coupon"
	"mymall/services/merchant-service/internal/svc"
)

func ClaimCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewClaimCouponLogic(r.Context(), svcCtx)
		l.ClaimCoupon(w, r)
	}
}

func ListMyCouponsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewListMyCouponsLogic(r.Context(), svcCtx)
		l.ListMyCoupons(w, r)
	}
}
