package product

import (
	"context"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductListLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductListLogic {
	return &GetProductListLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetProductListLogic) GetProductList(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	page, pageSize := req.Page, req.PageSize
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	shopID, _ := strconv.ParseUint("" /* was query:shop_id */, 10, 64)
	categoryID, _ := strconv.ParseUint("" /* was query:category_id */, 10, 64)
	orderBy := "" /* was query:order_by */
	data, err := plogic.NewCatalogLogic(l.svcCtx).GetProductListFiltered(ctx, pageReq, shopID, "on_sale", categoryID, orderBy)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return &types.PageListResp{List: data}, nil
}
