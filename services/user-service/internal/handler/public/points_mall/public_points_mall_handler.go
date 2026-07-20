package points_mall

import (
	"net/http"

	"mymall/services/user-service/internal/logic/public/points_mall"
	"mymall/services/user-service/internal/svc"
)

func ServePointsMallUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := points_mall.NewServePointsMallUploadLogic(r.Context(), svcCtx)
		l.ServePointsMallUpload(w, r)
	}
}
