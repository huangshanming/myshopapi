package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantListThemeSlotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMerchantListThemeSlotsLogic(r.Context(), svcCtx)
		l.MerchantListThemeSlots(w, r)
	}
}
