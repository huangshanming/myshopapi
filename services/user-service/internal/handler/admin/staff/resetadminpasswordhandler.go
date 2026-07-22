// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package staff

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/admin/staff"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func ResetAdminPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminResetPwdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := staff.NewResetAdminPasswordLogic(r.Context(), svcCtx)
		resp, err := l.ResetAdminPassword(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
