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

type MerchantUploadArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUploadArticleLogic(svcCtx *svc.ServiceContext) *MerchantUploadArticleLogic {
	return &MerchantUploadArticleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUploadArticleLogic) MerchantUploadArticle(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/article-uploads", nil, nil, req, hmerchant.NewArticleHandler(l.svcCtx).Upload)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
