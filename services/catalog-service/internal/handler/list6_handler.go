package handler

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/svc"
)

func List6Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewList6Logic(r.Context(), svcCtx)
		l.List6(w, r)
	}
}
