package points

import (
	"net/http"

	"mymall/services/user-service/internal/logic/user/points"
	"mymall/services/user-service/internal/svc"
)

func UserPointLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points.NewUserPointLogsLogic(r.Context(), svcCtx)
		l.UserPointLogs(w, r)
	}
}

func UserPointsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points.NewUserPointsLogic(r.Context(), svcCtx)
		l.UserPoints(w, r)
	}
}
