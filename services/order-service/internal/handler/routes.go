package handler

import (
	"net/http"

	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	"mymall/pkg/metrics"
	hadmin "mymall/services/order-service/internal/handler/admin"
	hmerchant "mymall/services/order-service/internal/handler/merchant"
	huser "mymall/services/order-service/internal/handler/user"
	svcMW "mymall/services/order-service/internal/middleware"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, svcCtx *svc.ServiceContext, healthReg *health.Registry, mws svcMW.Bundle) {
	orderUser := huser.NewOrderHandler(svcCtx)
	orderMerchant := hmerchant.NewOrderHandler(svcCtx)
	orderAdmin := hadmin.NewOrderHandler(svcCtx)
	reviewUser := huser.NewReviewHandler(svcCtx)
	reviewMerchant := hmerchant.NewReviewHandler(svcCtx)
	reviewAdmin := hadmin.NewReviewHandler(svcCtx)
	logisticsAdmin := hadmin.NewLogisticsHandler(svcCtx)

	server.AddRoutes(mws.Public([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: httpserver.Healthz("order-service")},
		{Method: http.MethodGet, Path: "/readyz", Handler: healthReg.ReadyHandler()},
		{Method: http.MethodGet, Path: "/metrics", Handler: metrics.Handler()},
		{Method: http.MethodGet, Path: "/api/v1/products/:id/reviews", Handler: reviewUser.ProductList},
	}))
	server.AddRoutes(mws.Authed([]rest.Route{
		{Method: http.MethodPost, Path: "/api/v1/orders", Handler: orderUser.Create},
		{Method: http.MethodPost, Path: "/api/v1/orders/coupon-preview", Handler: orderUser.CouponPreview},
		{Method: http.MethodGet, Path: "/api/v1/orders/status-counts", Handler: orderUser.StatusCounts},
		{Method: http.MethodGet, Path: "/api/v1/orders/after-sales", Handler: orderUser.UserAfterSales},
		{Method: http.MethodGet, Path: "/api/v1/orders", Handler: orderUser.List},
		{Method: http.MethodGet, Path: "/api/v1/orders/:id", Handler: orderUser.Detail},
		{Method: http.MethodPut, Path: "/api/v1/orders/:id/cancel", Handler: orderUser.Cancel},
		{Method: http.MethodPut, Path: "/api/v1/orders/:id/confirm-receive", Handler: orderUser.ConfirmReceive},
		{Method: http.MethodPost, Path: "/api/v1/orders/:id/after-sales", Handler: orderUser.CreateAfterSale},
		{Method: http.MethodGet, Path: "/api/v1/orders/:id/review-eligible", Handler: reviewUser.Eligible},
		{Method: http.MethodPost, Path: "/api/v1/orders/:id/reviews", Handler: reviewUser.Create},
		{Method: http.MethodGet, Path: "/api/v1/orders/:id/review", Handler: reviewUser.GetByOrder},
		{Method: http.MethodPost, Path: "/api/v1/user/review-uploads", Handler: reviewUser.Upload},
	}))
	server.AddRoutes(mws.MerchantOwner([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/merchant/reviews", Handler: reviewMerchant.MerchantList},
		{Method: http.MethodPut, Path: "/api/v1/merchant/reviews/:id/reply", Handler: reviewMerchant.MerchantReply},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/reviews/:id", Handler: reviewMerchant.MerchantDelete},
		{Method: http.MethodGet, Path: "/api/v1/merchant/orders", Handler: orderMerchant.MerchantList},
		{Method: http.MethodGet, Path: "/api/v1/merchant/orders/:id", Handler: orderMerchant.MerchantDetail},
		{Method: http.MethodPut, Path: "/api/v1/merchant/orders/:id/ship", Handler: orderMerchant.MerchantShip},
		{Method: http.MethodPut, Path: "/api/v1/merchant/orders/:id/complete", Handler: orderMerchant.MerchantComplete},
		{Method: http.MethodPut, Path: "/api/v1/merchant/orders/:id/remark", Handler: orderMerchant.MerchantRemark},
		{Method: http.MethodGet, Path: "/api/v1/merchant/after-sales", Handler: orderMerchant.MerchantAfterSales},
		{Method: http.MethodPut, Path: "/api/v1/merchant/after-sales/:id/handle", Handler: orderMerchant.MerchantHandleAfterSale},
	}))
	server.AddRoutes(mws.PlatformOrMerchant([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/logistics/options", Handler: logisticsAdmin.Options},
	}))
	server.AddRoutes(mws.PlatformAdmin([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/admin/reviews", Handler: reviewAdmin.AdminList},
		{Method: http.MethodDelete, Path: "/api/v1/admin/reviews/:id", Handler: reviewAdmin.AdminDelete},
		{Method: http.MethodGet, Path: "/api/v1/admin/orders", Handler: orderAdmin.AdminList},
		{Method: http.MethodGet, Path: "/api/v1/admin/orders/:id", Handler: orderAdmin.AdminDetail},
		{Method: http.MethodPut, Path: "/api/v1/admin/orders/:id/ship", Handler: orderAdmin.AdminShip},
		{Method: http.MethodPut, Path: "/api/v1/admin/orders/:id/complete", Handler: orderAdmin.AdminComplete},
		{Method: http.MethodPut, Path: "/api/v1/admin/orders/:id/remark", Handler: orderAdmin.AdminRemark},
		{Method: http.MethodGet, Path: "/api/v1/admin/after-sales", Handler: orderAdmin.AdminAfterSales},
		{Method: http.MethodPut, Path: "/api/v1/admin/after-sales/:id/handle", Handler: orderAdmin.AdminHandleAfterSale},
		{Method: http.MethodGet, Path: "/api/v1/admin/logistics", Handler: logisticsAdmin.AdminList},
		{Method: http.MethodPost, Path: "/api/v1/admin/logistics", Handler: logisticsAdmin.Create},
		{Method: http.MethodPut, Path: "/api/v1/admin/logistics/:id", Handler: logisticsAdmin.Update},
		{Method: http.MethodPut, Path: "/api/v1/admin/logistics/:id/status", Handler: logisticsAdmin.UpdateStatus},
		{Method: http.MethodDelete, Path: "/api/v1/admin/logistics/:id", Handler: logisticsAdmin.Delete},
	}))
}
