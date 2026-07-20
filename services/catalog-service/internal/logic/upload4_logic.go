package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type Upload4Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpload4Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Upload4Logic {
	return &Upload4Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Upload4Logic) Upload4(w http.ResponseWriter, r *http.Request) {
	padmin.NewShopUploadHandler().Upload(w, r)
}
