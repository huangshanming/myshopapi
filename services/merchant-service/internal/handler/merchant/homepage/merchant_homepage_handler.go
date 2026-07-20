package homepage

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/merchant/homepage"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantBuySlotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewMerchantBuySlotLogic(r.Context(), svcCtx)
		l.MerchantBuySlot(w, r)
	}
}

func MerchantListSlotOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewMerchantListSlotOrdersLogic(r.Context(), svcCtx)
		l.MerchantListSlotOrders(w, r)
	}
}

func MerchantListSlotPackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := homepage.NewMerchantListSlotPackagesLogic(r.Context(), svcCtx)
		l.MerchantListSlotPackages(w, r)
	}
}
