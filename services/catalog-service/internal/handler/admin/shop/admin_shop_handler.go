package shop

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/admin/shop"
	"mymall/services/catalog-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminUploadShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := shop.NewAdminUploadShopLogic(r.Context(), svcCtx)
		resp, err := l.AdminUploadShop(r.Context(), r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
