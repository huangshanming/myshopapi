package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cmerchant "mymall/services/catalog-service/internal/content/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type MerchantPatchArticleCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantPatchArticleCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantPatchArticleCommentLogic {
	return &MerchantPatchArticleCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantPatchArticleCommentLogic) MerchantPatchArticleComment(w http.ResponseWriter, r *http.Request) {
	cmerchant.NewArticleHandler(l.svcCtx).CommentPatch(w, r)
}
