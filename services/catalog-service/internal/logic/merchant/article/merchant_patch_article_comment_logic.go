package article

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hmerchant "mymall/services/catalog-service/internal/content/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantPatchArticleCommentLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantPatchArticleCommentLogic(svcCtx *svc.ServiceContext) *MerchantPatchArticleCommentLogic {
	return &MerchantPatchArticleCommentLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantPatchArticleCommentLogic) MerchantPatchArticleComment(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "PATCH", "/api/v1/merchant/article-comments/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewArticleHandler(l.svcCtx).CommentPatch)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
