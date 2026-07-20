package review

import (
	"net/http"

	"mymall/services/order-service/internal/logic/user/review"
	"mymall/services/order-service/internal/svc"
)

func EligibleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewEligibleLogic(r.Context(), svcCtx)
		l.Eligible(w, r)
	}
}

func GetByOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewGetByOrderLogic(r.Context(), svcCtx)
		l.GetByOrder(w, r)
	}
}

func UserCreateReviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewUserCreateReviewLogic(r.Context(), svcCtx)
		l.UserCreateReview(w, r)
	}
}

func UserUploadReviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := review.NewUserUploadReviewLogic(r.Context(), svcCtx)
		l.UserUploadReview(w, r)
	}
}
