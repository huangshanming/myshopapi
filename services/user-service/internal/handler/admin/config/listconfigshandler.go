// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/user-service/internal/logic/admin/config"
	"mymall/services/user-service/internal/svc"
)

func ListConfigsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := config.NewListConfigsLogic(r.Context(), svcCtx)
		resp, err := l.ListConfigs()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
