package order

import (
	"net/http"

	"mymall/services/order-service/internal/logic/admin/order"
	"mymall/services/order-service/internal/svc"
)

func AdminAfterSalesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewAdminAfterSalesLogic(r.Context(), svcCtx)
		l.AdminAfterSales(w, r)
	}
}

func AdminCompleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewAdminCompleteLogic(r.Context(), svcCtx)
		l.AdminComplete(w, r)
	}
}

func AdminDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewAdminDetailLogic(r.Context(), svcCtx)
		l.AdminDetail(w, r)
	}
}

func AdminHandleAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewAdminHandleAfterSaleLogic(r.Context(), svcCtx)
		l.AdminHandleAfterSale(w, r)
	}
}

func AdminListOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewAdminListOrdersLogic(r.Context(), svcCtx)
		l.AdminListOrders(w, r)
	}
}

func AdminRemarkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewAdminRemarkLogic(r.Context(), svcCtx)
		l.AdminRemark(w, r)
	}
}

func AdminShipHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewAdminShipLogic(r.Context(), svcCtx)
		l.AdminShip(w, r)
	}
}
