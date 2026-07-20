package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func UnreadNotificationCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUnreadNotificationCountLogic(r.Context(), svcCtx)
		l.UnreadNotificationCount(w, r)
	}
}
