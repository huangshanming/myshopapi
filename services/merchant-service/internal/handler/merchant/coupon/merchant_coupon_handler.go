package coupon

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/merchant/coupon"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantCopyCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewMerchantCopyCouponLogic(r.Context(), svcCtx)
		l.MerchantCopyCoupon(w, r)
	}
}

func MerchantCouponClaimsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewMerchantCouponClaimsLogic(r.Context(), svcCtx)
		l.MerchantCouponClaims(w, r)
	}
}

func MerchantCouponRedeemsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewMerchantCouponRedeemsLogic(r.Context(), svcCtx)
		l.MerchantCouponRedeems(w, r)
	}
}

func MerchantCouponStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewMerchantCouponStatsLogic(r.Context(), svcCtx)
		l.MerchantCouponStats(w, r)
	}
}

func MerchantCreateCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewMerchantCreateCouponLogic(r.Context(), svcCtx)
		l.MerchantCreateCoupon(w, r)
	}
}

func MerchantGrantCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewMerchantGrantCouponLogic(r.Context(), svcCtx)
		l.MerchantGrantCoupon(w, r)
	}
}

func MerchantListCouponsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewMerchantListCouponsLogic(r.Context(), svcCtx)
		l.MerchantListCoupons(w, r)
	}
}

func MerchantOffCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewMerchantOffCouponLogic(r.Context(), svcCtx)
		l.MerchantOffCoupon(w, r)
	}
}

func MerchantUpdateCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewMerchantUpdateCouponLogic(r.Context(), svcCtx)
		l.MerchantUpdateCoupon(w, r)
	}
}
