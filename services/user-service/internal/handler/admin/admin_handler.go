package admin

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin"
	"mymall/services/user-service/internal/svc"
)

func AdminAdjustWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminAdjustWalletLogic(r.Context(), svcCtx)
		l.AdminAdjustWallet(w, r)
	}
}

func AdminGetWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminGetWalletLogic(r.Context(), svcCtx)
		l.AdminGetWallet(w, r)
	}
}

func AdminListNotificationRecipientsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListNotificationRecipientsLogic(r.Context(), svcCtx)
		l.AdminListNotificationRecipients(w, r)
	}
}

func AdminListNotificationSendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListNotificationSendsLogic(r.Context(), svcCtx)
		l.AdminListNotificationSends(w, r)
	}
}

func AdminListTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListTasksLogic(r.Context(), svcCtx)
		l.AdminListTasks(w, r)
	}
}

func AdminListUserAddressesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListUserAddressesLogic(r.Context(), svcCtx)
		l.AdminListUserAddresses(w, r)
	}
}

func AdminSendNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminSendNotificationLogic(r.Context(), svcCtx)
		l.AdminSendNotification(w, r)
	}
}

func AdminUpdateTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminUpdateTaskLogic(r.Context(), svcCtx)
		l.AdminUpdateTask(w, r)
	}
}

func AdminWalletLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminWalletLogsLogic(r.Context(), svcCtx)
		l.AdminWalletLogs(w, r)
	}
}

func AssignAdminRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAssignAdminRolesLogic(r.Context(), svcCtx)
		l.AssignAdminRoles(w, r)
	}
}

func AssignRoleMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAssignRoleMenusLogic(r.Context(), svcCtx)
		l.AssignRoleMenus(w, r)
	}
}

func AuthMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAuthMeLogic(r.Context(), svcCtx)
		l.AuthMe(w, r)
	}
}

func CreateAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewCreateAdminLogic(r.Context(), svcCtx)
		l.CreateAdmin(w, r)
	}
}

func CreateMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewCreateMenuLogic(r.Context(), svcCtx)
		l.CreateMenu(w, r)
	}
}

func CreateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewCreateRoleLogic(r.Context(), svcCtx)
		l.CreateRole(w, r)
	}
}

func DeleteMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewDeleteMenuLogic(r.Context(), svcCtx)
		l.DeleteMenu(w, r)
	}
}

func DeleteRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewDeleteRoleLogic(r.Context(), svcCtx)
		l.DeleteRole(w, r)
	}
}

func GenerateUserTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewGenerateUserTokenLogic(r.Context(), svcCtx)
		l.GenerateUserToken(w, r)
	}
}

func GetAdminRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewGetAdminRolesLogic(r.Context(), svcCtx)
		l.GetAdminRoles(w, r)
	}
}

func GetRoleMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewGetRoleMenusLogic(r.Context(), svcCtx)
		l.GetRoleMenus(w, r)
	}
}

func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewGetUserLogic(r.Context(), svcCtx)
		l.GetUser(w, r)
	}
}

func ListAdminsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewListAdminsLogic(r.Context(), svcCtx)
		l.ListAdmins(w, r)
	}
}

func ListConfigsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewListConfigsLogic(r.Context(), svcCtx)
		l.ListConfigs(w, r)
	}
}

func ListRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewListRolesLogic(r.Context(), svcCtx)
		l.ListRoles(w, r)
	}
}

func ListUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewListUsersLogic(r.Context(), svcCtx)
		l.ListUsers(w, r)
	}
}

func MenuTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewMenuTreeLogic(r.Context(), svcCtx)
		l.MenuTree(w, r)
	}
}

func ResetAdminPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewResetAdminPasswordLogic(r.Context(), svcCtx)
		l.ResetAdminPassword(w, r)
	}
}

func ResetUserPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewResetUserPasswordLogic(r.Context(), svcCtx)
		l.ResetUserPassword(w, r)
	}
}

func SaveConfigsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewSaveConfigsLogic(r.Context(), svcCtx)
		l.SaveConfigs(w, r)
	}
}

func SetUserStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewSetUserStatusLogic(r.Context(), svcCtx)
		l.SetUserStatus(w, r)
	}
}

func UpdateMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewUpdateMenuLogic(r.Context(), svcCtx)
		l.UpdateMenu(w, r)
	}
}

func UpdateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewUpdateRoleLogic(r.Context(), svcCtx)
		l.UpdateRole(w, r)
	}
}

func UpdateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewUpdateUserLogic(r.Context(), svcCtx)
		l.UpdateUser(w, r)
	}
}
