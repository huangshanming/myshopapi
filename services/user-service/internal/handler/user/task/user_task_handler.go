package task

import (
	"net/http"

	"mymall/services/user-service/internal/logic/user/task"
	"mymall/services/user-service/internal/svc"
)

func UserCheckinHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := task.NewUserCheckinLogic(r.Context(), svcCtx)
		l.UserCheckin(w, r)
	}
}

func UserClaimHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := task.NewUserClaimLogic(r.Context(), svcCtx)
		l.UserClaim(w, r)
	}
}

func UserListTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := task.NewUserListTasksLogic(r.Context(), svcCtx)
		l.UserListTasks(w, r)
	}
}

func UserReportEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := task.NewUserReportEventLogic(r.Context(), svcCtx)
		l.UserReportEvent(w, r)
	}
}
