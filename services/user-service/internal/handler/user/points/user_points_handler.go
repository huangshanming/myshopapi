package points

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/services/user-service/internal/logic/user/points"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func UserPointLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := points.NewUserPointLogsLogic(svcCtx)
		resp, err := l.UserPointLogs(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func UserPointsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points.NewUserPointsLogic(svcCtx)
		resp, err := l.UserPoints(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
