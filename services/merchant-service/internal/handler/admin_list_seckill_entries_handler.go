package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func AdminListSeckillEntriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminListSeckillEntriesLogic(r.Context(), svcCtx)
		l.AdminListSeckillEntries(w, r)
	}
}
