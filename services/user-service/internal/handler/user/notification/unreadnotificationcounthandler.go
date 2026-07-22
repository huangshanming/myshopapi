// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/user/notification"
	"mymall/services/user-service/internal/svc"
)

func UnreadNotificationCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewUnreadNotificationCountLogic(r.Context(), svcCtx)
		resp, err := l.UnreadNotificationCount()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
