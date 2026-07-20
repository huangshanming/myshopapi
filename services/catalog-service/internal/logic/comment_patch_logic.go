package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cmerchant "mymall/services/catalog-service/internal/content/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type CommentPatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentPatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentPatchLogic {
	return &CommentPatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentPatchLogic) CommentPatch(w http.ResponseWriter, r *http.Request) {
	cmerchant.NewArticleHandler(l.svcCtx).CommentPatch(w, r)
}
