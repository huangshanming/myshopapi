// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/user/task"
	"mymall/services/user-service/internal/svc"
)

func UserListTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := task.NewUserListTasksLogic(r.Context(), svcCtx)
		resp, err := l.UserListTasks()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
