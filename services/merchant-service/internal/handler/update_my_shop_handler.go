package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func UpdateMyShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUpdateMyShopLogic(r.Context(), svcCtx)
		l.UpdateMyShop(w, r)
	}
}
