// Code generated manually for lottery-service.
package handler

import (
	"net/http"

	adminlottery "mymall/services/lottery-service/internal/handler/admin/lottery"
	publichealth "mymall/services/lottery-service/internal/handler/public/health"
	userlottery "mymall/services/lottery-service/internal/handler/user/lottery"
	"mymall/services/lottery-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{Method: http.MethodGet, Path: "/healthz", Handler: publichealth.HealthzHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/readyz", Handler: publichealth.ReadyzHandler(serverCtx)},
		},
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.RequestID},
			[]rest.Route{
				{Method: http.MethodGet, Path: "/api/v1/lottery/activity", Handler: userlottery.GetActivityHandler(serverCtx)},
			}...,
		),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.RequestID, serverCtx.GatewayIdentity},
			[]rest.Route{
				{Method: http.MethodPost, Path: "/api/v1/lottery/draw", Handler: userlottery.DrawHandler(serverCtx)},
				{Method: http.MethodGet, Path: "/api/v1/lottery/records", Handler: userlottery.ListRecordsHandler(serverCtx)},
				{Method: http.MethodPost, Path: "/api/v1/lottery/records/:id/address", Handler: userlottery.ClaimAddressHandler(serverCtx)},
			}...,
		),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.RequestID, serverCtx.GatewayIdentity, serverCtx.RequirePlatformAdmin},
			[]rest.Route{
				{Method: http.MethodGet, Path: "/api/v1/admin/lottery/activities", Handler: adminlottery.ListActivitiesHandler(serverCtx)},
				{Method: http.MethodPost, Path: "/api/v1/admin/lottery/activities", Handler: adminlottery.CreateActivityHandler(serverCtx)},
				{Method: http.MethodGet, Path: "/api/v1/admin/lottery/activities/:id", Handler: adminlottery.GetActivityHandler(serverCtx)},
				{Method: http.MethodPut, Path: "/api/v1/admin/lottery/activities/:id", Handler: adminlottery.UpdateActivityHandler(serverCtx)},
				{Method: http.MethodPut, Path: "/api/v1/admin/lottery/activities/:id/prizes", Handler: adminlottery.SavePrizesHandler(serverCtx)},
				{Method: http.MethodGet, Path: "/api/v1/admin/lottery/records", Handler: adminlottery.ListRecordsHandler(serverCtx)},
				{Method: http.MethodGet, Path: "/api/v1/admin/lottery/orders", Handler: adminlottery.ListOrdersHandler(serverCtx)},
				{Method: http.MethodPost, Path: "/api/v1/admin/lottery/orders/:id/ship", Handler: adminlottery.ShipOrderHandler(serverCtx)},
				{Method: http.MethodPost, Path: "/api/v1/admin/lottery/upload", Handler: adminlottery.UploadHandler(serverCtx)},
			}...,
		),
	)
}
