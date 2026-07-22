// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package points_mall

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/admin/points_mall"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func UpdatePointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PointsProductUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := points_mall.NewUpdatePointsProductLogic(r.Context(), svcCtx)
		resp, err := l.UpdatePointsProduct(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
