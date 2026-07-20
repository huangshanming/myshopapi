package merchant

import (
	"net/http"

	"mymall/services/order-service/internal/logic/merchant"
	"mymall/services/order-service/internal/svc"
)

func MerchantAfterSalesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantAfterSalesLogic(r.Context(), svcCtx)
		l.MerchantAfterSales(w, r)
	}
}

func MerchantCompleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantCompleteLogic(r.Context(), svcCtx)
		l.MerchantComplete(w, r)
	}
}

func MerchantDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantDeleteLogic(r.Context(), svcCtx)
		l.MerchantDelete(w, r)
	}
}

func MerchantDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantDetailLogic(r.Context(), svcCtx)
		l.MerchantDetail(w, r)
	}
}

func MerchantHandleAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantHandleAfterSaleLogic(r.Context(), svcCtx)
		l.MerchantHandleAfterSale(w, r)
	}
}

func MerchantListOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantListOrdersLogic(r.Context(), svcCtx)
		l.MerchantListOrders(w, r)
	}
}

func MerchantListReviewsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantListReviewsLogic(r.Context(), svcCtx)
		l.MerchantListReviews(w, r)
	}
}

func MerchantRemarkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantRemarkLogic(r.Context(), svcCtx)
		l.MerchantRemark(w, r)
	}
}

func MerchantReplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantReplyLogic(r.Context(), svcCtx)
		l.MerchantReply(w, r)
	}
}

func MerchantShipHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantShipLogic(r.Context(), svcCtx)
		l.MerchantShip(w, r)
	}
}
