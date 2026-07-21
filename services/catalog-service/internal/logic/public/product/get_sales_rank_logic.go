package product

import (
	"context"
	"mymall/pkg/appinput"
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

func (l *GetSalesRankLogic) GetSalesRank(ctx context.Context) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{}

	page, pageSize := in.Page()
	data, err := plogic.NewCatalogLogic(l.svcCtx).GetSalesRank(ctx, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return &types.AnyResp{Data: data}, nil
}
