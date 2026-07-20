package shop

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/admin/shop"
	"mymall/services/merchant-service/internal/svc"
)

func AdminCreateShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewAdminCreateShopLogic(r.Context(), svcCtx)
		l.AdminCreateShop(w, r)
	}
}

func AdminDisableShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewAdminDisableShopLogic(r.Context(), svcCtx)
		l.AdminDisableShop(w, r)
	}
}

func AdminEnableShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewAdminEnableShopLogic(r.Context(), svcCtx)
		l.AdminEnableShop(w, r)
	}
}

func AdminGetShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewAdminGetShopLogic(r.Context(), svcCtx)
		l.AdminGetShop(w, r)
	}
}

func AdminListShopsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewAdminListShopsLogic(r.Context(), svcCtx)
		l.AdminListShops(w, r)
	}
}

func AdminResetOwnerPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewAdminResetOwnerPasswordLogic(r.Context(), svcCtx)
		l.AdminResetOwnerPassword(w, r)
	}
}

func AdminUpdateShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewAdminUpdateShopLogic(r.Context(), svcCtx)
		l.AdminUpdateShop(w, r)
	}
}
