package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type CategoryCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryCreateLogic {
	return &CategoryCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryCreateLogic) CategoryCreate(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).CategoryCreate(w, r)
}
