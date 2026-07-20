package address

import (
	"net/http"

	"mymall/services/user-service/internal/logic/internalapi/address"
	"mymall/services/user-service/internal/svc"
)

func InternalGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewInternalGetLogic(r.Context(), svcCtx)
		l.InternalGet(w, r)
	}
}
