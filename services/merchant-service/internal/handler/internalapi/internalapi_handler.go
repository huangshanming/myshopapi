package internalapi

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/internalapi"
	"mymall/services/merchant-service/internal/svc"
)

func InternalLockCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalLockCouponLogic(r.Context(), svcCtx)
		l.InternalLockCoupon(w, r)
	}
}

func InternalMatchCouponsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalMatchCouponsLogic(r.Context(), svcCtx)
		l.InternalMatchCoupons(w, r)
	}
}

func InternalOrderGiftHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalOrderGiftLogic(r.Context(), svcCtx)
		l.InternalOrderGift(w, r)
	}
}

func InternalRedeemCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalRedeemCouponLogic(r.Context(), svcCtx)
		l.InternalRedeemCoupon(w, r)
	}
}

func InternalReturnCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalReturnCouponLogic(r.Context(), svcCtx)
		l.InternalReturnCoupon(w, r)
	}
}

func InternalUnlockCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalUnlockCouponLogic(r.Context(), svcCtx)
		l.InternalUnlockCoupon(w, r)
	}
}

func SeckillConsumeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewSeckillConsumeLogic(r.Context(), svcCtx)
		l.SeckillConsume(w, r)
	}
}

func SeckillRestoreHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewSeckillRestoreLogic(r.Context(), svcCtx)
		l.SeckillRestore(w, r)
	}
}
