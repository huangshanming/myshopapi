package logistics

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/order-service/internal/logic/shared/logistics"
	"mymall/services/order-service/internal/svc"
)

func LogisticsOptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logistics.NewLogisticsOptionsLogic(svcCtx)
		resp, err := l.LogisticsOptions(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
