package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func MarkNotificationReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMarkNotificationReadLogic(r.Context(), svcCtx)
		l.MarkNotificationRead(w, r)
	}
}
