package application

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/admin/application"
	"mymall/services/merchant-service/internal/svc"
)

func AdminApproveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := application.NewAdminApproveLogic(r.Context(), svcCtx)
		l.AdminApprove(w, r)
	}
}

func AdminListApplicationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := application.NewAdminListApplicationsLogic(r.Context(), svcCtx)
		l.AdminListApplications(w, r)
	}
}

func AdminRejectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := application.NewAdminRejectLogic(r.Context(), svcCtx)
		l.AdminReject(w, r)
	}
}
