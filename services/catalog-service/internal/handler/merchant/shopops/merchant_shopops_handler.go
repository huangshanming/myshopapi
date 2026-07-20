package shopops

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/merchant/shopops"
	"mymall/services/catalog-service/internal/svc"
)

func BindStaffHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shopops.NewBindStaffLogic(r.Context(), svcCtx)
		l.BindStaff(w, r)
	}
}

func ListMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shopops.NewListMenusLogic(r.Context(), svcCtx)
		l.ListMenus(w, r)
	}
}

func ListRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shopops.NewListRolesLogic(r.Context(), svcCtx)
		l.ListRoles(w, r)
	}
}

func ListStaffHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shopops.NewListStaffLogic(r.Context(), svcCtx)
		l.ListStaff(w, r)
	}
}

func MerchantAuthMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shopops.NewMerchantAuthMeLogic(r.Context(), svcCtx)
		l.MerchantAuthMe(w, r)
	}
}

func MerchantCreateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shopops.NewMerchantCreateRoleLogic(r.Context(), svcCtx)
		l.MerchantCreateRole(w, r)
	}
}

func MerchantUpdateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shopops.NewMerchantUpdateRoleLogic(r.Context(), svcCtx)
		l.MerchantUpdateRole(w, r)
	}
}

func RoleMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shopops.NewRoleMenusLogic(r.Context(), svcCtx)
		l.RoleMenus(w, r)
	}
}
