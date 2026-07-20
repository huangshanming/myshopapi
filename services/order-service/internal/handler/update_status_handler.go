package handler

import (
	"net/http"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

func UpdateStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUpdateStatusLogic(r.Context(), svcCtx)
		l.UpdateStatus(w, r)
	}
}
