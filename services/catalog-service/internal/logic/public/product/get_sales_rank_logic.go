package product

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hpublic "mymall/services/catalog-service/internal/product/app/public"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSalesRankLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetSalesRankLogic(svcCtx *svc.ServiceContext) *GetSalesRankLogic {
	return &GetSalesRankLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *GetSalesRankLogic) GetSalesRank(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/products/sales-rank", nil, nil, nil, hpublic.NewCatalogHandler(l.svcCtx).GetSalesRank)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
