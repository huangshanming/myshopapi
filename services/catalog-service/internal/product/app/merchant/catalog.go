package merchant

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"
)

func (h *CatalogHandler) MerchantListProducts(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少 shop_id")
	}
	page, pageSize := in.Page()
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	status := in.QueryGet("status")
	data, err := h.logic.GetProductListFiltered(ctx, pageReq, shopID, status, 0, "")
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *CatalogHandler) MerchantCreateProduct(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少 shop_id")
	}
	var req types.MerchantProductReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if req.Name == "" || req.SalePrice == 0 || req.CategoryID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	p, err := h.logic.CreateProduct(ctx, shopID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *CatalogHandler) MerchantUpdateProduct(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少 shop_id")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	var req types.MerchantProductReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if req.Name == "" || req.SalePrice == 0 || req.CategoryID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateProductByShop(ctx, id, shopID, req); err != nil {
		return nil, xerr.New(http.StatusForbidden, err.Error())
	}
	return nil, nil
}

func (h *CatalogHandler) MerchantSetStatus(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少 shop_id")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	var body types.SetStatusReq
	if err := appinput.BindBody(in, &body); err != nil || body.Status == "" {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.SetProductStatus(ctx, id, shopID, body.Status); err != nil {
		return nil, xerr.New(http.StatusForbidden, err.Error())
	}
	return nil, nil
}
