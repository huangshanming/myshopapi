package admin

import (
	"net/http"

	"mymall/services/order-service/internal/logic/admin"
	"mymall/services/order-service/internal/svc"
)

func AdminAfterSalesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminAfterSalesLogic(r.Context(), svcCtx)
		l.AdminAfterSales(w, r)
	}
}

func AdminCompleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCompleteLogic(r.Context(), svcCtx)
		l.AdminComplete(w, r)
	}
}

func AdminCreateLogisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCreateLogisticsLogic(r.Context(), svcCtx)
		l.AdminCreateLogistics(w, r)
	}
}

func AdminDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminDeleteLogic(r.Context(), svcCtx)
		l.AdminDelete(w, r)
	}
}

func AdminDeleteLogisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminDeleteLogisticsLogic(r.Context(), svcCtx)
		l.AdminDeleteLogistics(w, r)
	}
}

func AdminDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminDetailLogic(r.Context(), svcCtx)
		l.AdminDetail(w, r)
	}
}

func AdminHandleAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminHandleAfterSaleLogic(r.Context(), svcCtx)
		l.AdminHandleAfterSale(w, r)
	}
}

func AdminListLogisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListLogisticsLogic(r.Context(), svcCtx)
		l.AdminListLogistics(w, r)
	}
}

func AdminListOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListOrdersLogic(r.Context(), svcCtx)
		l.AdminListOrders(w, r)
	}
}

func AdminListReviewsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListReviewsLogic(r.Context(), svcCtx)
		l.AdminListReviews(w, r)
	}
}

func AdminRemarkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminRemarkLogic(r.Context(), svcCtx)
		l.AdminRemark(w, r)
	}
}

func AdminShipHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminShipLogic(r.Context(), svcCtx)
		l.AdminShip(w, r)
	}
}

func AdminUpdateLogisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminUpdateLogisticsLogic(r.Context(), svcCtx)
		l.AdminUpdateLogistics(w, r)
	}
}

func UpdateStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewUpdateStatusLogic(r.Context(), svcCtx)
		l.UpdateStatus(w, r)
	}
}
