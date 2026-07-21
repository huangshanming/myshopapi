package address

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/internalapi/address"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func InternalGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InternalAddressReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := address.NewInternalGetLogic(svcCtx)
		resp, err := l.InternalGet(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
