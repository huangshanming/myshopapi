package notification

import (
	"net/http"

	"mymall/services/user-service/internal/logic/user/notification"
	"mymall/services/user-service/internal/svc"
)

func ListNotificationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewListNotificationsLogic(r.Context(), svcCtx)
		l.ListNotifications(w, r)
	}
}

func MarkAllNotificationsReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewMarkAllNotificationsReadLogic(r.Context(), svcCtx)
		l.MarkAllNotificationsRead(w, r)
	}
}

func MarkNotificationReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewMarkNotificationReadLogic(r.Context(), svcCtx)
		l.MarkNotificationRead(w, r)
	}
}

func UnreadNotificationCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewUnreadNotificationCountLogic(r.Context(), svcCtx)
		l.UnreadNotificationCount(w, r)
	}
}
