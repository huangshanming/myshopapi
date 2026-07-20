package handler

import (
	"net/http"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

func AdminList3Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminList3Logic(r.Context(), svcCtx)
		l.AdminList3(w, r)
	}
}
