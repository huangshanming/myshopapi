package public

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/httpserver"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/uploadpath"
)

type ServePointsMallUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServePointsMallUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServePointsMallUploadLogic {
	return &ServePointsMallUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServePointsMallUploadLogic) ServePointsMallUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("points-mall", httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
