package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func InternalCreateNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewInternalCreateNotificationLogic(r.Context(), svcCtx)
		l.InternalCreateNotification(w, r)
	}
}
