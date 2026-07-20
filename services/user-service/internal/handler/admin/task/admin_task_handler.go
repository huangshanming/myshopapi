package task

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin/task"
	"mymall/services/user-service/internal/svc"
)

func AdminListTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := task.NewAdminListTasksLogic(r.Context(), svcCtx)
		l.AdminListTasks(w, r)
	}
}

func AdminUpdateTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := task.NewAdminUpdateTaskLogic(r.Context(), svcCtx)
		l.AdminUpdateTask(w, r)
	}
}
