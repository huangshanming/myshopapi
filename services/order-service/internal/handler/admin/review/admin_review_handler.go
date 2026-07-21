package review

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/services/order-service/internal/logic/admin/review"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"
)

func AdminDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := review.NewAdminDeleteLogic(svcCtx)
		resp, err := l.AdminDelete(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func AdminListReviewsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := review.NewAdminListReviewsLogic(svcCtx)
		resp, err := l.AdminListReviews(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
