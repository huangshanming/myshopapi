package user

import (
	"net/http"

	"mymall/services/order-service/internal/logic/user"
	"mymall/services/order-service/internal/svc"
)

func CancelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewCancelLogic(r.Context(), svcCtx)
		l.Cancel(w, r)
	}
}

func ConfirmReceiveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewConfirmReceiveLogic(r.Context(), svcCtx)
		l.ConfirmReceive(w, r)
	}
}

func CouponPreviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewCouponPreviewLogic(r.Context(), svcCtx)
		l.CouponPreview(w, r)
	}
}

func CreateAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewCreateAfterSaleLogic(r.Context(), svcCtx)
		l.CreateAfterSale(w, r)
	}
}

func EligibleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewEligibleLogic(r.Context(), svcCtx)
		l.Eligible(w, r)
	}
}

func GetByOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetByOrderLogic(r.Context(), svcCtx)
		l.GetByOrder(w, r)
	}
}

func StatusCountsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewStatusCountsLogic(r.Context(), svcCtx)
		l.StatusCounts(w, r)
	}
}

func UserAfterSalesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserAfterSalesLogic(r.Context(), svcCtx)
		l.UserAfterSales(w, r)
	}
}

func UserCreateOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserCreateOrderLogic(r.Context(), svcCtx)
		l.UserCreateOrder(w, r)
	}
}

func UserCreateReviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserCreateReviewLogic(r.Context(), svcCtx)
		l.UserCreateReview(w, r)
	}
}

func UserGetOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserGetOrderLogic(r.Context(), svcCtx)
		l.UserGetOrder(w, r)
	}
}

func UserListOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserListOrdersLogic(r.Context(), svcCtx)
		l.UserListOrders(w, r)
	}
}

func UserUploadReviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserUploadReviewLogic(r.Context(), svcCtx)
		l.UserUploadReview(w, r)
	}
}
