package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"
)

func (h *CatalogHandler) AdminListProducts(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	status := in.QueryGet("status")
	data, err := h.logic.GetProductListFiltered(ctx, pageReq, shopID, status, 0, "")
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *CatalogHandler) AdminForceOffSale(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	if err := h.logic.ForceOffSale(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CatalogHandler) AdminListCategories(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.ListAllCategories(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *CatalogHandler) AdminCreateCategory(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.CategoryReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if req.Name == "" {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	cat, err := h.logic.CreateCategory(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return cat, nil
}

func (h *CatalogHandler) AdminUpdateCategory(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "分类ID无效")
	}
	var req types.CategoryReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if req.Name == "" {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateCategory(ctx, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CatalogHandler) AdminDeleteCategory(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "分类ID无效")
	}
	if err := h.logic.DeleteCategory(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
