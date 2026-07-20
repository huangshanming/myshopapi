package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type Detail4Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetail4Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Detail4Logic {
	return &Detail4Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Detail4Logic) Detail4(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).Detail(w, r)
}
