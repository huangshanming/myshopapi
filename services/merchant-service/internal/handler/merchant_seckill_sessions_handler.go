package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantSeckillSessionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMerchantSeckillSessionsLogic(r.Context(), svcCtx)
		l.MerchantSeckillSessions(w, r)
	}
}
