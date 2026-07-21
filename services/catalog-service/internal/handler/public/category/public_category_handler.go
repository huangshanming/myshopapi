package category

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/services/catalog-service/internal/logic/public/category"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"
)

func GetCategoryDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewGetCategoryDetailLogic(svcCtx)
		resp, err := l.GetCategoryDetail(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func GetCategoryListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := category.NewGetCategoryListLogic(svcCtx)
		resp, err := l.GetCategoryList(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
