package shop

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/public/shop"
	"mymall/services/merchant-service/internal/svc"
)

func PublicGetShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewPublicGetShopLogic(r.Context(), svcCtx)
		l.PublicGetShop(w, r)
	}
}

func PublicHomeSlotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewPublicHomeSlotsLogic(r.Context(), svcCtx)
		l.PublicHomeSlots(w, r)
	}
}

func PublicListShopsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewPublicListShopsLogic(r.Context(), svcCtx)
		l.PublicListShops(w, r)
	}
}

func PublicThemeTilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewPublicThemeTilesLogic(r.Context(), svcCtx)
		l.PublicThemeTiles(w, r)
	}
}
