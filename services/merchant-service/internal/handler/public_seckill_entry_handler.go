package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func PublicSeckillEntryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewPublicSeckillEntryLogic(r.Context(), svcCtx)
		l.PublicSeckillEntry(w, r)
	}
}
