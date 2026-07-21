package task

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/internalapi/task"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func InternalEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TaskEventReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewInternalEventLogic(svcCtx)
		err := l.InternalEvent(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
