package role

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin/role"
	"mymall/services/user-service/internal/svc"
)

func AssignRoleMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewAssignRoleMenusLogic(r.Context(), svcCtx)
		l.AssignRoleMenus(w, r)
	}
}

func CreateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewCreateRoleLogic(r.Context(), svcCtx)
		l.CreateRole(w, r)
	}
}

func DeleteRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewDeleteRoleLogic(r.Context(), svcCtx)
		l.DeleteRole(w, r)
	}
}

func GetRoleMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewGetRoleMenusLogic(r.Context(), svcCtx)
		l.GetRoleMenus(w, r)
	}
}

func ListRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewListRolesLogic(r.Context(), svcCtx)
		l.ListRoles(w, r)
	}
}

func UpdateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewUpdateRoleLogic(r.Context(), svcCtx)
		l.UpdateRole(w, r)
	}
}
