package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hpublic "mymall/services/user-service/internal/httpapi/public"
	"mymall/services/user-service/internal/svc"
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
	hpublic.ServePointsMallUpload(w, r)
}
