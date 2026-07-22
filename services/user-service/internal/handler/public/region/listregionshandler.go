// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package region

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/public/region"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func ListRegionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegionListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := region.NewListRegionsLogic(r.Context(), svcCtx)
		resp, err := l.ListRegions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
