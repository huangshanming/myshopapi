package review

import (
	"net/http"

	"mymall/services/order-service/internal/logic/admin/review"
	"mymall/services/order-service/internal/svc"
)

func AdminDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewAdminDeleteLogic(r.Context(), svcCtx)
		l.AdminDelete(w, r)
	}
}

func AdminListReviewsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewAdminListReviewsLogic(r.Context(), svcCtx)
		l.AdminListReviews(w, r)
	}
}
