package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListProductsLogic {
	return &MerchantListProductsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListProductsLogic) MerchantListProducts(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := req.Page, req.PageSize
	catID, _ := strconv.ParseUint("" /* was query:category_id */, 10, 64)
	f := repository.ProductListFilter{
		ShopID: shopID, Name: "" /* was query:name */, ProductNo: "" /* was query:product_no */,
		CategoryID: catID, Status: "" /* was query:status */, ProductType: "" /* was query:product_type */,
		StockWarnOnly: "" /* was query:stock_warn */ == "1",
		Page:          page, PageSize: pageSize, OrderBy: "" /* was query:order_by */,
		Recycle: "" /* was query:recycle */ == "1",
	}
	data, err := plogic.NewProductAdminLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
