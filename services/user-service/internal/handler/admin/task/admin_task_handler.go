package task

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/services/user-service/internal/logic/admin/task"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func AdminListTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := task.NewAdminListTasksLogic(svcCtx)
		resp, err := l.AdminListTasks(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func AdminUpdateTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewAdminUpdateTaskLogic(svcCtx)
		resp, err := l.AdminUpdateTask(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
