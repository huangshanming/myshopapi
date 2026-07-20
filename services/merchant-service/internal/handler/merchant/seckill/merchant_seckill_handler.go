package seckill

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/merchant/seckill"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantApplySeckillHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewMerchantApplySeckillLogic(r.Context(), svcCtx)
		l.MerchantApplySeckill(w, r)
	}
}

func MerchantListSeckillEntriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewMerchantListSeckillEntriesLogic(r.Context(), svcCtx)
		l.MerchantListSeckillEntries(w, r)
	}
}

func MerchantSeckillSessionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewMerchantSeckillSessionsLogic(r.Context(), svcCtx)
		l.MerchantSeckillSessions(w, r)
	}
}

func MerchantSetSeckillAutoRenewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewMerchantSetSeckillAutoRenewLogic(r.Context(), svcCtx)
		l.MerchantSetSeckillAutoRenew(w, r)
	}
}
