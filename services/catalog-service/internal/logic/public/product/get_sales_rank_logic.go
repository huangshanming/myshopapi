package product

import (
	"context"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

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

func (l *GetSalesRankLogic) GetSalesRank(ctx context.Context, req *types.PageReq) (resp *types.SalesRankResp, err error) {
	data, err := plogic.NewCatalogLogic(l.svcCtx).GetSalesRank(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return &types.SalesRankResp{Data: data}, nil
}
