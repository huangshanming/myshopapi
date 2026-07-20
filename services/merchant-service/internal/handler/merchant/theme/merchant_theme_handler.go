package theme

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/merchant/theme"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantBuyThemeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewMerchantBuyThemeLogic(r.Context(), svcCtx)
		l.MerchantBuyTheme(w, r)
	}
}

func MerchantListThemeOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewMerchantListThemeOrdersLogic(r.Context(), svcCtx)
		l.MerchantListThemeOrders(w, r)
	}
}

func MerchantListThemePackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewMerchantListThemePackagesLogic(r.Context(), svcCtx)
		l.MerchantListThemePackages(w, r)
	}
}

func MerchantListThemeSlotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := theme.NewMerchantListThemeSlotsLogic(r.Context(), svcCtx)
		l.MerchantListThemeSlots(w, r)
	}
}
