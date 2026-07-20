package homepage

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/admin/homepage"
	"mymall/services/merchant-service/internal/svc"
)

func AdminCreateSlotPackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewAdminCreateSlotPackageLogic(r.Context(), svcCtx)
		l.AdminCreateSlotPackage(w, r)
	}
}

func AdminGrantSlotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewAdminGrantSlotLogic(r.Context(), svcCtx)
		l.AdminGrantSlot(w, r)
	}
}

func AdminListSlotOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewAdminListSlotOrdersLogic(r.Context(), svcCtx)
		l.AdminListSlotOrders(w, r)
	}
}

func AdminListSlotPackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewAdminListSlotPackagesLogic(r.Context(), svcCtx)
		l.AdminListSlotPackages(w, r)
	}
}

func AdminListSlotSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewAdminListSlotSettingsLogic(r.Context(), svcCtx)
		l.AdminListSlotSettings(w, r)
	}
}

func AdminUpdateSlotPackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewAdminUpdateSlotPackageLogic(r.Context(), svcCtx)
		l.AdminUpdateSlotPackage(w, r)
	}
}

func AdminUpdateSlotSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewAdminUpdateSlotSettingsLogic(r.Context(), svcCtx)
		l.AdminUpdateSlotSettings(w, r)
	}
}
