package staff

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin/staff"
	"mymall/services/user-service/internal/svc"
)

func AssignAdminRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := staff.NewAssignAdminRolesLogic(r.Context(), svcCtx)
		l.AssignAdminRoles(w, r)
	}
}

func CreateAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := staff.NewCreateAdminLogic(r.Context(), svcCtx)
		l.CreateAdmin(w, r)
	}
}

func GetAdminRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := staff.NewGetAdminRolesLogic(r.Context(), svcCtx)
		l.GetAdminRoles(w, r)
	}
}

func ListAdminsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := staff.NewListAdminsLogic(r.Context(), svcCtx)
		l.ListAdmins(w, r)
	}
}

func ResetAdminPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := staff.NewResetAdminPasswordLogic(r.Context(), svcCtx)
		l.ResetAdminPassword(w, r)
	}
}
