package config

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin/config"
	"mymall/services/user-service/internal/svc"
)

func ListConfigsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := config.NewListConfigsLogic(r.Context(), svcCtx)
		l.ListConfigs(w, r)
	}
}

func SaveConfigsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := config.NewSaveConfigsLogic(r.Context(), svcCtx)
		l.SaveConfigs(w, r)
	}
}
