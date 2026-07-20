package review

import (
	"net/http"

	"mymall/services/order-service/internal/logic/public/review"
	"mymall/services/order-service/internal/svc"
)

func PublicListProductReviewsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewPublicListProductReviewsLogic(r.Context(), svcCtx)
		l.PublicListProductReviews(w, r)
	}
}
