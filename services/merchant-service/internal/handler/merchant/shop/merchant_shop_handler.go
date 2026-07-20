package shop

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/merchant/shop"
	"mymall/services/merchant-service/internal/svc"
)

func ApplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewApplyLogic(r.Context(), svcCtx)
		l.Apply(w, r)
	}
}

func MyShopsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewMyShopsLogic(r.Context(), svcCtx)
		l.MyShops(w, r)
	}
}

func UpdateMyShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewUpdateMyShopLogic(r.Context(), svcCtx)
		l.UpdateMyShop(w, r)
	}
}
