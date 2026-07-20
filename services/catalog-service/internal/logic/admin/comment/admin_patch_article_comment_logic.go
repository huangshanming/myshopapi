package comment

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminPatchArticleCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminPatchArticleCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminPatchArticleCommentLogic {
	return &AdminPatchArticleCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminPatchArticleCommentLogic) AdminPatchArticleComment(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).CommentPatch(w, r)
}
