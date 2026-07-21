package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hmerchant "mymall/services/catalog-service/internal/content/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantPatchArticleCommentLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantPatchArticleCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantPatchArticleCommentLogic {
	return &MerchantPatchArticleCommentLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantPatchArticleCommentLogic) MerchantPatchArticleComment(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hmerchant.NewArticleHandler(l.svcCtx).CommentPatch(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
