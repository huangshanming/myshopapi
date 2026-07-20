package logistics

import (
	"net/http"

	"mymall/services/order-service/internal/logic/admin/logistics"
	"mymall/services/order-service/internal/svc"
)

func AdminCreateLogisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logistics.NewAdminCreateLogisticsLogic(r.Context(), svcCtx)
		l.AdminCreateLogistics(w, r)
	}
}

func AdminDeleteLogisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logistics.NewAdminDeleteLogisticsLogic(r.Context(), svcCtx)
		l.AdminDeleteLogistics(w, r)
	}
}

func AdminListLogisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logistics.NewAdminListLogisticsLogic(r.Context(), svcCtx)
		l.AdminListLogistics(w, r)
	}
}

func AdminUpdateLogisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logistics.NewAdminUpdateLogisticsLogic(r.Context(), svcCtx)
		l.AdminUpdateLogistics(w, r)
	}
}

func UpdateStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logistics.NewUpdateStatusLogic(r.Context(), svcCtx)
		l.UpdateStatus(w, r)
	}
}
