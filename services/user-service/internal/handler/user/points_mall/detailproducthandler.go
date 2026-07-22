// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package points_mall

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/user/points_mall"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func DetailProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := points_mall.NewDetailProductLogic(r.Context(), svcCtx)
		resp, err := l.DetailProduct(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
