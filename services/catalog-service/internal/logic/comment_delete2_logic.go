package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type CommentDelete2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentDelete2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentDelete2Logic {
	return &CommentDelete2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentDelete2Logic) CommentDelete2(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).CommentDelete(w, r)
}
