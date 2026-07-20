package seckill

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/internalapi/seckill"
	"mymall/services/merchant-service/internal/svc"
)

func SeckillConsumeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewSeckillConsumeLogic(r.Context(), svcCtx)
		l.SeckillConsume(w, r)
	}
}

func SeckillRestoreHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewSeckillRestoreLogic(r.Context(), svcCtx)
		l.SeckillRestore(w, r)
	}
}
