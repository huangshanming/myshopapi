package user

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/user"
	"mymall/services/merchant-service/internal/svc"
)

func ClaimCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewClaimCouponLogic(r.Context(), svcCtx)
		l.ClaimCoupon(w, r)
	}
}

func DetailUserPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewDetailUserPointsOrderLogic(r.Context(), svcCtx)
		l.DetailUserPointsOrder(w, r)
	}
}

func ExchangeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewExchangeLogic(r.Context(), svcCtx)
		l.Exchange(w, r)
	}
}

func ListMyCouponsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewListMyCouponsLogic(r.Context(), svcCtx)
		l.ListMyCoupons(w, r)
	}
}

func ListUserPointsOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewListUserPointsOrdersLogic(r.Context(), svcCtx)
		l.ListUserPointsOrders(w, r)
	}
}
