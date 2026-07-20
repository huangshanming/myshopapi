package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cmerchant "mymall/services/catalog-service/internal/content/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type List4Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList4Logic(ctx context.Context, svcCtx *svc.ServiceContext) *List4Logic {
	return &List4Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List4Logic) List4(w http.ResponseWriter, r *http.Request) {
	cmerchant.NewArticleHandler(l.svcCtx).List(w, r)
}
