package product

import (
	"context"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductListLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductListLogic {
	return &GetProductListLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *GetProductListLogic) GetProductList(ctx context.Context, req *types.PublicProductListReq) (resp *types.PageListResp, err error) {
	pageReq := &pagination.PageReq{Page: req.Page, PageSize: req.PageSize}
	data, err := plogic.NewCatalogLogic(l.svcCtx).GetProductListFiltered(ctx, pageReq, req.ShopId, "on_sale", req.CategoryId, req.OrderBy)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return types.FromPaged(data), nil
}
