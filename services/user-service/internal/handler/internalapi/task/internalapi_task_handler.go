package task

import (
	"net/http"

	"mymall/services/user-service/internal/logic/internalapi/task"
	"mymall/services/user-service/internal/svc"
)

func InternalEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := task.NewInternalEventLogic(r.Context(), svcCtx)
		l.InternalEvent(w, r)
	}
}
