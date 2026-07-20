package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func AdminListNotificationRecipientsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminListNotificationRecipientsLogic(r.Context(), svcCtx)
		l.AdminListNotificationRecipients(w, r)
	}
}
