package profile

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/user/profile"
	"mymall/services/user-service/internal/svc"
)

func UserProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := profile.NewUserProfileLogic(svcCtx)
		resp, err := l.UserProfile(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
