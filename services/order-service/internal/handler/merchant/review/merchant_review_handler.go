package review

import (
	"net/http"

	"mymall/services/order-service/internal/logic/merchant/review"
	"mymall/services/order-service/internal/svc"
)

func MerchantDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewMerchantDeleteLogic(r.Context(), svcCtx)
		l.MerchantDelete(w, r)
	}
}

func MerchantListReviewsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewMerchantListReviewsLogic(r.Context(), svcCtx)
		l.MerchantListReviews(w, r)
	}
}

func MerchantReplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewMerchantReplyLogic(r.Context(), svcCtx)
		l.MerchantReply(w, r)
	}
}
