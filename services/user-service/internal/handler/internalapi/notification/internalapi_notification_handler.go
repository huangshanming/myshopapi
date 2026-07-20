package notification

import (
	"net/http"

	"mymall/services/user-service/internal/logic/internalapi/notification"
	"mymall/services/user-service/internal/svc"
)

func InternalCreateNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewInternalCreateNotificationLogic(r.Context(), svcCtx)
		l.InternalCreateNotification(w, r)
	}
}
