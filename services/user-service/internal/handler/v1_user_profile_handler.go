package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func V1UserProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewV1UserProfileLogic(r.Context(), svcCtx)
		l.V1UserProfile(w, r)
	}
}
