package points_mall

import (
	"net/http"

	"mymall/services/user-service/internal/uploadpath"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ServePointsMallUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FilePathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		http.ServeFile(w, r, uploadpath.Abs("points-mall", req.File))
	}
}
