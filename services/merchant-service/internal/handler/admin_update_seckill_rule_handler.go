package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func AdminUpdateSeckillRuleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminUpdateSeckillRuleLogic(r.Context(), svcCtx)
		l.AdminUpdateSeckillRule(w, r)
	}
}
