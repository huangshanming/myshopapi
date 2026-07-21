package notification

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/internalapi/notification"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func InternalCreateNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NotifyCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := notification.NewInternalCreateNotificationLogic(svcCtx)
		resp, err := l.InternalCreateNotification(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
