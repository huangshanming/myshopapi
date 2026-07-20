package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type CategoryList2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryList2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryList2Logic {
	return &CategoryList2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryList2Logic) CategoryList2(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).CategoryList(w, r)
}
