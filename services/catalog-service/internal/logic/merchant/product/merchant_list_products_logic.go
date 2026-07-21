package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListProductsLogic {
	return &MerchantListProductsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantListProductsLogic) MerchantListProducts(ctx context.Context, req *types.MerchantProductListReq) (resp *types.PageListResp, err error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	f := repository.ProductListFilter{
		ShopID: shopID, Name: req.Name, ProductNo: req.ProductNo,
		CategoryID: req.CategoryId, Status: req.Status, ProductType: req.ProductType,
		StockWarnOnly: req.StockWarn == "1",
		Page: req.Page, PageSize: req.PageSize, OrderBy: req.OrderBy,
		Recycle: req.Recycle == "1",
	}
	data, err := plogic.NewProductAdminLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
