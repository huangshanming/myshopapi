// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package points_mall

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/admin/points_mall"
	"mymall/services/user-service/internal/svc"
)

func UploadPointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewUploadPointsProductLogic(r.Context(), svcCtx)
		resp, err := l.UploadPointsProduct()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
