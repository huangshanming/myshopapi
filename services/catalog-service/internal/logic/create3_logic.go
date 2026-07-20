package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type Create3Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreate3Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Create3Logic {
	return &Create3Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Create3Logic) Create3(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).Create(w, r)
}
