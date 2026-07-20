package notification

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin/notification"
	"mymall/services/user-service/internal/svc"
)

func AdminListNotificationRecipientsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewAdminListNotificationRecipientsLogic(r.Context(), svcCtx)
		l.AdminListNotificationRecipients(w, r)
	}
}

func AdminListNotificationSendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewAdminListNotificationSendsLogic(r.Context(), svcCtx)
		l.AdminListNotificationSends(w, r)
	}
}

func AdminSendNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewAdminSendNotificationLogic(r.Context(), svcCtx)
		l.AdminSendNotification(w, r)
	}
}
