package order

import (
	"net/http"

	"mymall/services/order-service/internal/logic/user/order"
	"mymall/services/order-service/internal/svc"
)

func CancelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewCancelLogic(r.Context(), svcCtx)
		l.Cancel(w, r)
	}
}

func ConfirmReceiveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewConfirmReceiveLogic(r.Context(), svcCtx)
		l.ConfirmReceive(w, r)
	}
}

func CouponPreviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewCouponPreviewLogic(r.Context(), svcCtx)
		l.CouponPreview(w, r)
	}
}

func CreateAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewCreateAfterSaleLogic(r.Context(), svcCtx)
		l.CreateAfterSale(w, r)
	}
}

func StatusCountsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewStatusCountsLogic(r.Context(), svcCtx)
		l.StatusCounts(w, r)
	}
}

func UserAfterSalesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewUserAfterSalesLogic(r.Context(), svcCtx)
		l.UserAfterSales(w, r)
	}
}

func UserCreateOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewUserCreateOrderLogic(r.Context(), svcCtx)
		l.UserCreateOrder(w, r)
	}
}

func UserGetOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewUserGetOrderLogic(r.Context(), svcCtx)
		l.UserGetOrder(w, r)
	}
}

func UserListOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewUserListOrdersLogic(r.Context(), svcCtx)
		l.UserListOrders(w, r)
	}
}
