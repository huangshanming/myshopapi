// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/user/address"
	"mymall/services/user-service/internal/svc"
)

func UserListAddressesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewUserListAddressesLogic(r.Context(), svcCtx)
		resp, err := l.UserListAddresses()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
