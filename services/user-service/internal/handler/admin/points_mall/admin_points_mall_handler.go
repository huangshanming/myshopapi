package points_mall

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin/points_mall"
	"mymall/services/user-service/internal/svc"
)

func CancelPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewCancelPointsOrderLogic(r.Context(), svcCtx)
		l.CancelPointsOrder(w, r)
	}
}

func CompletePointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewCompletePointsOrderLogic(r.Context(), svcCtx)
		l.CompletePointsOrder(w, r)
	}
}

func CreatePointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewCreatePointsProductLogic(r.Context(), svcCtx)
		l.CreatePointsProduct(w, r)
	}
}

func DeletePointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewDeletePointsProductLogic(r.Context(), svcCtx)
		l.DeletePointsProduct(w, r)
	}
}

func DetailPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewDetailPointsOrderLogic(r.Context(), svcCtx)
		l.DetailPointsOrder(w, r)
	}
}

func DetailPointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewDetailPointsProductLogic(r.Context(), svcCtx)
		l.DetailPointsProduct(w, r)
	}
}

func ListPointsOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewListPointsOrdersLogic(r.Context(), svcCtx)
		l.ListPointsOrders(w, r)
	}
}

func ListPointsProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewListPointsProductsLogic(r.Context(), svcCtx)
		l.ListPointsProducts(w, r)
	}
}

func RemarkPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewRemarkPointsOrderLogic(r.Context(), svcCtx)
		l.RemarkPointsOrder(w, r)
	}
}

func SetPointsProductStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewSetPointsProductStatusLogic(r.Context(), svcCtx)
		l.SetPointsProductStatus(w, r)
	}
}

func ShipPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewShipPointsOrderLogic(r.Context(), svcCtx)
		l.ShipPointsOrder(w, r)
	}
}

func UpdatePointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewUpdatePointsProductLogic(r.Context(), svcCtx)
		l.UpdatePointsProduct(w, r)
	}
}

func UploadPointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewUploadPointsProductLogic(r.Context(), svcCtx)
		l.UploadPointsProduct(w, r)
	}
}
