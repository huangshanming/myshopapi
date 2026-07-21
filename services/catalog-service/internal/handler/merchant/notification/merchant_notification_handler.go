package notification

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/services/catalog-service/internal/logic/merchant/notification"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"
)

func MerchantListNotificationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := notification.NewMerchantListNotificationsLogic(r.Context(), svcCtx)
		resp, err := l.MerchantListNotifications(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func MerchantMarkAllNotificationsReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewMerchantMarkAllNotificationsReadLogic(r.Context(), svcCtx)
		resp, err := l.MerchantMarkAllNotificationsRead(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func MerchantMarkNotificationReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := notification.NewMerchantMarkNotificationReadLogic(r.Context(), svcCtx)
		resp, err := l.MerchantMarkNotificationRead(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func MerchantUnreadNotificationCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewMerchantUnreadNotificationCountLogic(r.Context(), svcCtx)
		resp, err := l.MerchantUnreadNotificationCount(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
