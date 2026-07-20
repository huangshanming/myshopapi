package theme

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/admin/theme"
	"mymall/services/merchant-service/internal/svc"
)

func AdminCreateThemePackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewAdminCreateThemePackageLogic(r.Context(), svcCtx)
		l.AdminCreateThemePackage(w, r)
	}
}

func AdminGrantThemeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewAdminGrantThemeLogic(r.Context(), svcCtx)
		l.AdminGrantTheme(w, r)
	}
}

func AdminListThemeOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewAdminListThemeOrdersLogic(r.Context(), svcCtx)
		l.AdminListThemeOrders(w, r)
	}
}

func AdminListThemePackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewAdminListThemePackagesLogic(r.Context(), svcCtx)
		l.AdminListThemePackages(w, r)
	}
}

func AdminListThemeSlotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewAdminListThemeSlotsLogic(r.Context(), svcCtx)
		l.AdminListThemeSlots(w, r)
	}
}

func AdminUpdateThemePackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewAdminUpdateThemePackageLogic(r.Context(), svcCtx)
		l.AdminUpdateThemePackage(w, r)
	}
}

func AdminUpdateThemeSlotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewAdminUpdateThemeSlotLogic(r.Context(), svcCtx)
		l.AdminUpdateThemeSlot(w, r)
	}
}
