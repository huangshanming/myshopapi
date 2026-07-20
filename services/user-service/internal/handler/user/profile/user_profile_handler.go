package profile

import (
	"net/http"

	"mymall/services/user-service/internal/logic/user/profile"
	"mymall/services/user-service/internal/svc"
)

func UserProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := profile.NewUserProfileLogic(r.Context(), svcCtx)
		l.UserProfile(w, r)
	}
}
