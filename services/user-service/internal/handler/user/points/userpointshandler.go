// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package points

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/user/points"
	"mymall/services/user-service/internal/svc"
)

func UserPointsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points.NewUserPointsLogic(r.Context(), svcCtx)
		resp, err := l.UserPoints()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
