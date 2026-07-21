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

type MerchantUpdateArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUpdateArticleLogic(svcCtx *svc.ServiceContext) *MerchantUpdateArticleLogic {
	return &MerchantUpdateArticleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUpdateArticleLogic) MerchantUpdateArticle(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/merchant/articles/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewArticleHandler(l.svcCtx).Update)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
