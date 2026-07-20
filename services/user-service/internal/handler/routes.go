package handler

import (
	"net/http"

	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	"mymall/pkg/metrics"
	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/handler/admin"
	hinternal "mymall/services/user-service/internal/handler/internalapi"
	hpublic "mymall/services/user-service/internal/handler/public"
	huser "mymall/services/user-service/internal/handler/user"
	svcMW "mymall/services/user-service/internal/middleware"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, svcCtx *svc.ServiceContext, healthReg *health.Registry, mws svcMW.Bundle) {
	userHandler := huser.NewUserHandler(svcCtx)
	adminHandler := hadmin.NewAdminHandler(svcCtx)
	walletUser := huser.NewWalletHandler(svcCtx)
	walletAdmin := hadmin.NewWalletHandler(svcCtx)
	walletInternal := hinternal.NewWalletHandler(svcCtx)
	addressUser := huser.NewAddressHandler(svcCtx)
	addressAdmin := hadmin.NewAddressHandler(svcCtx)
	addressInternal := hinternal.NewAddressHandler(svcCtx)
	regionHandler := hpublic.NewRegionHandler(svcCtx)
	taskUser := huser.NewTaskHandler(svcCtx)
	taskAdmin := hadmin.NewTaskHandler(svcCtx)
	taskInternal := hinternal.NewTaskHandler(svcCtx)
	notifInternal := hinternal.NewNotificationHandler(svcCtx)

	withPerm := func(code string, h http.HandlerFunc) http.HandlerFunc {
		if code == "" {
			return h
		}
		return pkgmw.RequirePermission(adminHandler, code)(h)
	}

	server.AddRoutes(mws.Public([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: httpserver.Healthz("user-service")},
		{Method: http.MethodGet, Path: "/readyz", Handler: healthReg.ReadyHandler()},
		{Method: http.MethodGet, Path: "/metrics", Handler: metrics.Handler()},
		{Method: http.MethodPost, Path: "/api/v1/user/login", Handler: userHandler.Login},
		{Method: http.MethodPost, Path: "/api/v1/user/register", Handler: userHandler.Register},
		{Method: http.MethodPost, Path: "/api/v1/user/wallet/freeze", Handler: walletInternal.Freeze},
		{Method: http.MethodPost, Path: "/api/v1/user/wallet/unfreeze", Handler: walletInternal.Unfreeze},
		{Method: http.MethodPost, Path: "/api/v1/user/wallet/settle", Handler: walletInternal.Settle},
		{Method: http.MethodGet, Path: "/api/v1/user/addresses/internal", Handler: addressInternal.InternalGet},
		{Method: http.MethodPost, Path: "/api/v1/internal/notifications", Handler: notifInternal.InternalCreateNotification},
		{Method: http.MethodPost, Path: "/api/v1/internal/tasks/events", Handler: taskInternal.InternalEvent},
		{Method: http.MethodPost, Path: "/api/v1/internal/points/deduct", Handler: taskInternal.InternalDeductPoints},
		{Method: http.MethodPost, Path: "/api/v1/internal/points/refund", Handler: taskInternal.InternalRefundPoints},
		{Method: http.MethodGet, Path: "/api/v1/regions", Handler: regionHandler.List},
		{Method: http.MethodGet, Path: "/api/v1/regions/tree", Handler: regionHandler.Tree},
	}))
	server.AddRoutes(mws.UserFlexible([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/user/profile", Handler: userHandler.Profile},
		{Method: http.MethodGet, Path: "/api/v1/user/wallet", Handler: walletUser.UserGetWallet},
		{Method: http.MethodGet, Path: "/api/v1/user/wallet/logs", Handler: walletUser.UserWalletLogs},
		{Method: http.MethodGet, Path: "/api/v1/user/addresses", Handler: addressUser.List},
		{Method: http.MethodPost, Path: "/api/v1/user/addresses", Handler: addressUser.Create},
		{Method: http.MethodPut, Path: "/api/v1/user/addresses/:id", Handler: addressUser.Update},
		{Method: http.MethodDelete, Path: "/api/v1/user/addresses/:id", Handler: addressUser.Delete},
		{Method: http.MethodPut, Path: "/api/v1/user/addresses/:id/default", Handler: addressUser.SetDefault},
		{Method: http.MethodGet, Path: "/api/v1/user/notifications", Handler: userHandler.ListNotifications},
		{Method: http.MethodGet, Path: "/api/v1/user/notifications/unread-count", Handler: userHandler.UnreadNotificationCount},
		{Method: http.MethodPost, Path: "/api/v1/user/notifications/:id/read", Handler: userHandler.MarkNotificationRead},
		{Method: http.MethodPost, Path: "/api/v1/user/notifications/read-all", Handler: userHandler.MarkAllNotificationsRead},
		{Method: http.MethodGet, Path: "/api/v1/user/points", Handler: taskUser.UserPoints},
		{Method: http.MethodGet, Path: "/api/v1/user/points/logs", Handler: taskUser.UserPointLogs},
		{Method: http.MethodGet, Path: "/api/v1/user/tasks", Handler: taskUser.UserListTasks},
		{Method: http.MethodPost, Path: "/api/v1/user/tasks/checkin", Handler: taskUser.UserCheckin},
		{Method: http.MethodPost, Path: "/api/v1/user/tasks/:code/claim", Handler: taskUser.UserClaim},
		{Method: http.MethodPost, Path: "/api/v1/user/tasks/events", Handler: taskUser.UserReportEvent},
	}))
	server.AddRoutes(mws.PlatformAdmin([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/admin/auth/me", Handler: withPerm("", adminHandler.AuthMe)},
		{Method: http.MethodGet, Path: "/api/v1/admin/menus", Handler: withPerm("system:menu:list", adminHandler.MenuTree)},
		{Method: http.MethodPost, Path: "/api/v1/admin/menus", Handler: withPerm("system:menu:add", adminHandler.CreateMenu)},
		{Method: http.MethodPut, Path: "/api/v1/admin/menus/:id", Handler: withPerm("system:menu:edit", adminHandler.UpdateMenu)},
		{Method: http.MethodDelete, Path: "/api/v1/admin/menus/:id", Handler: withPerm("system:menu:delete", adminHandler.DeleteMenu)},
		{Method: http.MethodGet, Path: "/api/v1/admin/roles", Handler: withPerm("system:role:list", adminHandler.ListRoles)},
		{Method: http.MethodPost, Path: "/api/v1/admin/roles", Handler: withPerm("system:role:add", adminHandler.CreateRole)},
		{Method: http.MethodPut, Path: "/api/v1/admin/roles/:id", Handler: withPerm("system:role:edit", adminHandler.UpdateRole)},
		{Method: http.MethodDelete, Path: "/api/v1/admin/roles/:id", Handler: withPerm("system:role:delete", adminHandler.DeleteRole)},
		{Method: http.MethodGet, Path: "/api/v1/admin/roles/:id/menus", Handler: withPerm("system:role:list", adminHandler.GetRoleMenus)},
		{Method: http.MethodPut, Path: "/api/v1/admin/roles/:id/menus", Handler: withPerm("system:role:assign", adminHandler.AssignRoleMenus)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users", Handler: withPerm("system:user:list", adminHandler.ListUsers)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id", Handler: withPerm("system:user:list", adminHandler.GetUser)},
		{Method: http.MethodPut, Path: "/api/v1/admin/users/:id", Handler: withPerm("system:user:edit", adminHandler.UpdateUser)},
		{Method: http.MethodPut, Path: "/api/v1/admin/users/:id/status", Handler: withPerm("system:user:status", adminHandler.SetUserStatus)},
		{Method: http.MethodPut, Path: "/api/v1/admin/users/:id/password", Handler: withPerm("system:user:reset", adminHandler.ResetUserPassword)},
		{Method: http.MethodPost, Path: "/api/v1/admin/users/:id/token", Handler: withPerm("system:user:list", adminHandler.GenerateUserToken)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id/wallet", Handler: withPerm("system:user:wallet", walletAdmin.AdminGetWallet)},
		{Method: http.MethodPost, Path: "/api/v1/admin/users/:id/wallet/adjust", Handler: withPerm("system:user:wallet", walletAdmin.AdminAdjustWallet)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id/wallet/logs", Handler: withPerm("system:user:wallet", walletAdmin.AdminWalletLogs)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id/addresses", Handler: withPerm("system:user:list", addressAdmin.AdminList)},
		{Method: http.MethodGet, Path: "/api/v1/admin/admins", Handler: withPerm("system:admin:list", adminHandler.ListAdmins)},
		{Method: http.MethodPost, Path: "/api/v1/admin/admins", Handler: withPerm("system:admin:add", adminHandler.CreateAdmin)},
		{Method: http.MethodGet, Path: "/api/v1/admin/admins/:id/roles", Handler: withPerm("system:admin:list", adminHandler.GetAdminRoles)},
		{Method: http.MethodPut, Path: "/api/v1/admin/admins/:id/roles", Handler: withPerm("system:admin:assign", adminHandler.AssignAdminRoles)},
		{Method: http.MethodPut, Path: "/api/v1/admin/admins/:id/password", Handler: withPerm("system:admin:reset", adminHandler.ResetAdminPassword)},
		{Method: http.MethodGet, Path: "/api/v1/admin/configs", Handler: withPerm("system:config:list", adminHandler.ListConfigs)},
		{Method: http.MethodPut, Path: "/api/v1/admin/configs", Handler: withPerm("system:config:edit", adminHandler.SaveConfigs)},
		{Method: http.MethodPost, Path: "/api/v1/admin/notifications/send", Handler: withPerm("business:message:send", adminHandler.AdminSendNotification)},
		{Method: http.MethodGet, Path: "/api/v1/admin/notifications/sends", Handler: withPerm("business:message:send", adminHandler.AdminListNotificationSends)},
		{Method: http.MethodGet, Path: "/api/v1/admin/notifications/sends/:id/recipients", Handler: withPerm("business:message:send", adminHandler.AdminListNotificationRecipients)},
		{Method: http.MethodGet, Path: "/api/v1/admin/tasks", Handler: withPerm("marketing:task:list", taskAdmin.AdminList)},
		{Method: http.MethodPut, Path: "/api/v1/admin/tasks/:id", Handler: withPerm("marketing:task:edit", taskAdmin.AdminUpdate)},
	}))
}
