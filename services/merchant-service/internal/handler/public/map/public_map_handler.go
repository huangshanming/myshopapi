package mapapi

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	maplogic "mymall/services/merchant-service/internal/logic/public/map"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
)

func PublicMapGeocoderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MapGeocoderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := maplogic.NewPublicMapGeocoderLogic(r.Context(), svcCtx)
		resp, err := l.PublicMapGeocoder(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func PublicMapConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := maplogic.NewPublicMapConfigLogic(r.Context(), svcCtx)
		resp, err := l.PublicMapConfig(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
