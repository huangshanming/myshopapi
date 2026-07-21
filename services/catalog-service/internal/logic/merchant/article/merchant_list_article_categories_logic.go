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

type MerchantListArticleCategoriesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListArticleCategoriesLogic(svcCtx *svc.ServiceContext) *MerchantListArticleCategoriesLogic {
	return &MerchantListArticleCategoriesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListArticleCategoriesLogic) MerchantListArticleCategories(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/merchant/article-categories", nil, url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}, nil, hmerchant.NewArticleHandler(l.svcCtx).CategoryList)
	if err != nil {
		return nil, err
	}
	var out types.PageListResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		var list interface{}
		if err2 := httpinvoke.Decode(raw, &list); err2 == nil {
			return &types.PageListResp{List: list}, nil
		}
		return nil, err
	}
	return &out, nil
}
