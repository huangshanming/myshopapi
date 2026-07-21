package product

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hpublic "mymall/services/catalog-service/internal/product/app/public"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSalesRankLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetSalesRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSalesRankLogic {
	return &GetSalesRankLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetSalesRankLogic) GetSalesRank(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hpublic.NewCatalogHandler(l.svcCtx).GetSalesRank(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
