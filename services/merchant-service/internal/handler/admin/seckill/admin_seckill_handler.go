package seckill

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/admin/seckill"
	"mymall/services/merchant-service/internal/svc"
)

func AdminGetSeckillRuleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewAdminGetSeckillRuleLogic(r.Context(), svcCtx)
		l.AdminGetSeckillRule(w, r)
	}
}

func AdminListSeckillEntriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewAdminListSeckillEntriesLogic(r.Context(), svcCtx)
		l.AdminListSeckillEntries(w, r)
	}
}

func AdminListSeckillSessionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewAdminListSeckillSessionsLogic(r.Context(), svcCtx)
		l.AdminListSeckillSessions(w, r)
	}
}

func AdminUpdateSeckillRuleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seckill.NewAdminUpdateSeckillRuleLogic(r.Context(), svcCtx)
		l.AdminUpdateSeckillRule(w, r)
	}
}
