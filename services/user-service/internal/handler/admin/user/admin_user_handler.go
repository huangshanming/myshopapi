package user

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin/user"
	"mymall/services/user-service/internal/svc"
)

func AdminAdjustWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewAdminAdjustWalletLogic(r.Context(), svcCtx)
		l.AdminAdjustWallet(w, r)
	}
}

func AdminGetWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewAdminGetWalletLogic(r.Context(), svcCtx)
		l.AdminGetWallet(w, r)
	}
}

func AdminListUserAddressesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewAdminListUserAddressesLogic(r.Context(), svcCtx)
		l.AdminListUserAddresses(w, r)
	}
}

func AdminWalletLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewAdminWalletLogsLogic(r.Context(), svcCtx)
		l.AdminWalletLogs(w, r)
	}
}

func GenerateUserTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGenerateUserTokenLogic(r.Context(), svcCtx)
		l.GenerateUserToken(w, r)
	}
}

func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserLogic(r.Context(), svcCtx)
		l.GetUser(w, r)
	}
}

func ListUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewListUsersLogic(r.Context(), svcCtx)
		l.ListUsers(w, r)
	}
}

func ResetUserPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewResetUserPasswordLogic(r.Context(), svcCtx)
		l.ResetUserPassword(w, r)
	}
}

func SetUserStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewSetUserStatusLogic(r.Context(), svcCtx)
		l.SetUserStatus(w, r)
	}
}

func UpdateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUpdateUserLogic(r.Context(), svcCtx)
		l.UpdateUser(w, r)
	}
}
