package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func AdminListNotificationSendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminListNotificationSendsLogic(r.Context(), svcCtx)
		l.AdminListNotificationSends(w, r)
	}
}
