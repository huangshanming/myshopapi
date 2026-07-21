package points_mall

import (
	"net/http"

	"mymall/services/user-service/internal/logic/user/points_mall"
	"mymall/services/user-service/internal/svc"
)

func DetailUserPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewDetailUserPointsOrderLogic(r.Context(), svcCtx)
		l.DetailUserPointsOrder(w, r)
	}
}

func ExchangeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewExchangeLogic(r.Context(), svcCtx)
		l.Exchange(w, r)
	}
}

func ListUserPointsOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewListUserPointsOrdersLogic(r.Context(), svcCtx)
		l.ListUserPointsOrders(w, r)
	}
}

func ListPointsGoodsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewPointProductLogic(r.Context(), svcCtx)
		l.ListPointsProduct(w, r)
	}
}
