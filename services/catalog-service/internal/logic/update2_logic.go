package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cmerchant "mymall/services/catalog-service/internal/content/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type Update2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdate2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Update2Logic {
	return &Update2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Update2Logic) Update2(w http.ResponseWriter, r *http.Request) {
	cmerchant.NewArticleHandler(l.svcCtx).Update(w, r)
}
