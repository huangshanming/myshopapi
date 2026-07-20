package user

import (
	"net/http"

	"mymall/services/user-service/internal/logic/user"
	"mymall/services/user-service/internal/svc"
)

func ListNotificationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewListNotificationsLogic(r.Context(), svcCtx)
		l.ListNotifications(w, r)
	}
}

func MarkAllNotificationsReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewMarkAllNotificationsReadLogic(r.Context(), svcCtx)
		l.MarkAllNotificationsRead(w, r)
	}
}

func MarkNotificationReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewMarkNotificationReadLogic(r.Context(), svcCtx)
		l.MarkNotificationRead(w, r)
	}
}

func SetDefaultHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewSetDefaultLogic(r.Context(), svcCtx)
		l.SetDefault(w, r)
	}
}

func UnreadNotificationCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUnreadNotificationCountLogic(r.Context(), svcCtx)
		l.UnreadNotificationCount(w, r)
	}
}

func UserCheckinHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserCheckinLogic(r.Context(), svcCtx)
		l.UserCheckin(w, r)
	}
}

func UserClaimHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserClaimLogic(r.Context(), svcCtx)
		l.UserClaim(w, r)
	}
}

func UserCreateAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserCreateAddressLogic(r.Context(), svcCtx)
		l.UserCreateAddress(w, r)
	}
}

func UserDeleteAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserDeleteAddressLogic(r.Context(), svcCtx)
		l.UserDeleteAddress(w, r)
	}
}

func UserGetWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserGetWalletLogic(r.Context(), svcCtx)
		l.UserGetWallet(w, r)
	}
}

func UserListAddressesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserListAddressesLogic(r.Context(), svcCtx)
		l.UserListAddresses(w, r)
	}
}

func UserListTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserListTasksLogic(r.Context(), svcCtx)
		l.UserListTasks(w, r)
	}
}

func UserPointLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserPointLogsLogic(r.Context(), svcCtx)
		l.UserPointLogs(w, r)
	}
}

func UserPointsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserPointsLogic(r.Context(), svcCtx)
		l.UserPoints(w, r)
	}
}

func UserProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserProfileLogic(r.Context(), svcCtx)
		l.UserProfile(w, r)
	}
}

func UserReportEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserReportEventLogic(r.Context(), svcCtx)
		l.UserReportEvent(w, r)
	}
}

func UserUpdateAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserUpdateAddressLogic(r.Context(), svcCtx)
		l.UserUpdateAddress(w, r)
	}
}

func UserWalletLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserWalletLogsLogic(r.Context(), svcCtx)
		l.UserWalletLogs(w, r)
	}
}
