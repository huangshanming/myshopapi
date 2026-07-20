package seckill

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/public/seckill"
	"mymall/services/merchant-service/internal/svc"
)

func PublicSeckillCurrentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewPublicSeckillCurrentLogic(r.Context(), svcCtx)
		l.PublicSeckillCurrent(w, r)
	}
}

func PublicSeckillEntryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewPublicSeckillEntryLogic(r.Context(), svcCtx)
		l.PublicSeckillEntry(w, r)
	}
}

func PublicSeckillListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewPublicSeckillListLogic(r.Context(), svcCtx)
		l.PublicSeckillList(w, r)
	}
}
