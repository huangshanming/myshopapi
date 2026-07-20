package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func AdminListCouponsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminListCouponsLogic(r.Context(), svcCtx)
		l.AdminListCoupons(w, r)
	}
}
