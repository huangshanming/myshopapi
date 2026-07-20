package coupon

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/admin/coupon"
	"mymall/services/merchant-service/internal/svc"
)

func AdminCopyCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewAdminCopyCouponLogic(r.Context(), svcCtx)
		l.AdminCopyCoupon(w, r)
	}
}

func AdminCouponClaimsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewAdminCouponClaimsLogic(r.Context(), svcCtx)
		l.AdminCouponClaims(w, r)
	}
}

func AdminCouponRedeemsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewAdminCouponRedeemsLogic(r.Context(), svcCtx)
		l.AdminCouponRedeems(w, r)
	}
}

func AdminCouponStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewAdminCouponStatsLogic(r.Context(), svcCtx)
		l.AdminCouponStats(w, r)
	}
}

func AdminCreateCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewAdminCreateCouponLogic(r.Context(), svcCtx)
		l.AdminCreateCoupon(w, r)
	}
}

func AdminGrantCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewAdminGrantCouponLogic(r.Context(), svcCtx)
		l.AdminGrantCoupon(w, r)
	}
}

func AdminListCouponsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewAdminListCouponsLogic(r.Context(), svcCtx)
		l.AdminListCoupons(w, r)
	}
}

func AdminOffCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewAdminOffCouponLogic(r.Context(), svcCtx)
		l.AdminOffCoupon(w, r)
	}
}

func AdminUpdateCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewAdminUpdateCouponLogic(r.Context(), svcCtx)
		l.AdminUpdateCoupon(w, r)
	}
}
