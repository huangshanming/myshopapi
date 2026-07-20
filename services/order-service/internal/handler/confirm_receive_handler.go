package handler

import (
	"net/http"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

func ConfirmReceiveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewConfirmReceiveLogic(r.Context(), svcCtx)
		l.ConfirmReceive(w, r)
	}
}
