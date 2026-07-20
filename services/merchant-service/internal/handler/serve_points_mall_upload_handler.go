package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func ServePointsMallUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewServePointsMallUploadLogic(r.Context(), svcCtx)
		l.ServePointsMallUpload(w, r)
	}
}
