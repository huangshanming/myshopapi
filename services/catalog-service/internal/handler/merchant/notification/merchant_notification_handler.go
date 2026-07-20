package notification

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/merchant/notification"
	"mymall/services/catalog-service/internal/svc"
)

func MerchantListNotificationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewMerchantListNotificationsLogic(r.Context(), svcCtx)
		l.MerchantListNotifications(w, r)
	}
}

func MerchantMarkAllNotificationsReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewMerchantMarkAllNotificationsReadLogic(r.Context(), svcCtx)
		l.MerchantMarkAllNotificationsRead(w, r)
	}
}

func MerchantMarkNotificationReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewMerchantMarkNotificationReadLogic(r.Context(), svcCtx)
		l.MerchantMarkNotificationRead(w, r)
	}
}

func MerchantUnreadNotificationCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewMerchantUnreadNotificationCountLogic(r.Context(), svcCtx)
		l.MerchantUnreadNotificationCount(w, r)
	}
}
