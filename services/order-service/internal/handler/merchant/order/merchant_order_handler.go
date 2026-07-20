package order

import (
	"net/http"

	"mymall/services/order-service/internal/logic/merchant/order"
	"mymall/services/order-service/internal/svc"
)

func MerchantAfterSalesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewMerchantAfterSalesLogic(r.Context(), svcCtx)
		l.MerchantAfterSales(w, r)
	}
}

func MerchantCompleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewMerchantCompleteLogic(r.Context(), svcCtx)
		l.MerchantComplete(w, r)
	}
}

func MerchantDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewMerchantDetailLogic(r.Context(), svcCtx)
		l.MerchantDetail(w, r)
	}
}

func MerchantHandleAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewMerchantHandleAfterSaleLogic(r.Context(), svcCtx)
		l.MerchantHandleAfterSale(w, r)
	}
}

func MerchantListOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewMerchantListOrdersLogic(r.Context(), svcCtx)
		l.MerchantListOrders(w, r)
	}
}

func MerchantRemarkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewMerchantRemarkLogic(r.Context(), svcCtx)
		l.MerchantRemark(w, r)
	}
}

func MerchantShipHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := order.NewMerchantShipLogic(r.Context(), svcCtx)
		l.MerchantShip(w, r)
	}
}
