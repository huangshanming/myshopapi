package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func PublicSeckillCurrentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewPublicSeckillCurrentLogic(r.Context(), svcCtx)
		l.PublicSeckillCurrent(w, r)
	}
}
