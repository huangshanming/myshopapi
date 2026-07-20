package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type CategoryDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryDeleteLogic {
	return &CategoryDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryDeleteLogic) CategoryDelete(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).CategoryDelete(w, r)
}
