package coupon

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/internalapi/coupon"
	"mymall/services/merchant-service/internal/svc"
)

func InternalLockCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewInternalLockCouponLogic(r.Context(), svcCtx)
		l.InternalLockCoupon(w, r)
	}
}

func InternalMatchCouponsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewInternalMatchCouponsLogic(r.Context(), svcCtx)
		l.InternalMatchCoupons(w, r)
	}
}

func InternalOrderGiftHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewInternalOrderGiftLogic(r.Context(), svcCtx)
		l.InternalOrderGift(w, r)
	}
}

func InternalRedeemCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewInternalRedeemCouponLogic(r.Context(), svcCtx)
		l.InternalRedeemCoupon(w, r)
	}
}

func InternalReturnCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewInternalReturnCouponLogic(r.Context(), svcCtx)
		l.InternalReturnCoupon(w, r)
	}
}

func InternalUnlockCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewInternalUnlockCouponLogic(r.Context(), svcCtx)
		l.InternalUnlockCoupon(w, r)
	}
}
