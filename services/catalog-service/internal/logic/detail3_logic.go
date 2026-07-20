package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cmerchant "mymall/services/catalog-service/internal/content/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type Detail3Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetail3Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Detail3Logic {
	return &Detail3Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Detail3Logic) Detail3(w http.ResponseWriter, r *http.Request) {
	cmerchant.NewArticleHandler(l.svcCtx).Detail(w, r)
}
