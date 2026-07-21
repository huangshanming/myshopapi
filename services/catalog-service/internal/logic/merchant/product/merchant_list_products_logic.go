package product

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"net/http"
	"net/url"
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
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := in.Page()
	catID, _ := strconv.ParseUint(in.QueryGet("category_id"), 10, 64)
	f := repository.ProductListFilter{
		ShopID: shopID, Name: in.QueryGet("name"), ProductNo: in.QueryGet("product_no"),
		CategoryID: catID, Status: in.QueryGet("status"), ProductType: in.QueryGet("product_type"),
		StockWarnOnly: in.QueryGet("stock_warn") == "1",
		Page:          page, PageSize: pageSize, OrderBy: in.QueryGet("order_by"),
		Recycle: in.QueryGet("recycle") == "1",
	}
	data, err := plogic.NewProductAdminLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
