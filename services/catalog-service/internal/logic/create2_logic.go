package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cmerchant "mymall/services/catalog-service/internal/content/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type Create2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreate2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Create2Logic {
	return &Create2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Create2Logic) Create2(w http.ResponseWriter, r *http.Request) {
	cmerchant.NewArticleHandler(l.svcCtx).Create(w, r)
}
