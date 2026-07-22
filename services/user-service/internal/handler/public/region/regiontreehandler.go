// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package region

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/public/region"
	"mymall/services/user-service/internal/svc"
)

func RegionTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := region.NewRegionTreeLogic(r.Context(), svcCtx)
		resp, err := l.RegionTree()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
