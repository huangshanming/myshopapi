package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type Upload3Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpload3Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Upload3Logic {
	return &Upload3Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Upload3Logic) Upload3(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).Upload(w, r)
}
