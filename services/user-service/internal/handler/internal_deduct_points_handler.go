package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func InternalDeductPointsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewInternalDeductPointsLogic(r.Context(), svcCtx)
		l.InternalDeductPoints(w, r)
	}
}
