package public

import (
	"context"
	"mymall/pkg/apidoc/dto"
	"mymall/pkg/appinput"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

var (
	_ dto.ProductListResp
	_ dto.ProductDetail
	_ dto.CategoryListResp
	_ dto.CategoryDetailResp
)

// GetProductList 商品列表
// @Summary      商品列表
// @Description  分页查询在售商品
// @Tags         商品
// @Produce      json
// @Param        page       query  int  false  "页码"       default(1)
// @Param        page_size  query  int  false  "每页数量"   default(10)
// @Success      200  {object}  apidoc.Response{data=dto.ProductListResp}  "查询成功"
// @Router       /api/v1/products/list [get]
func (h *CatalogHandler) GetProductList(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	categoryID, _ := strconv.ParseUint(in.QueryGet("category_id"), 10, 64)
	orderBy := in.QueryGet("order_by")
	data, err := h.logic.GetProductListFiltered(ctx, pageReq, shopID, "on_sale", categoryID, orderBy)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return data, nil
}

// GetSalesRank 今日必买销量榜
func (h *CatalogHandler) GetSalesRank(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	data, err := h.logic.GetSalesRank(ctx, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return data, nil
}

// GetProductDetail 商品详情
// @Summary      商品详情
// @Description  根据商品 ID 查询详情
// @Tags         商品
// @Produce      json
// @Param        id  query  int  true  "商品 ID"
// @Success      200  {object}  apidoc.Response{data=dto.ProductDetail}  "查询成功"
// @Router       /api/v1/products/detail [get]
func (h *CatalogHandler) GetProductDetail(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.QueryGet("id"), 10, 64)
	if id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	data, err := h.logic.GetProductDetail(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerr.New(http.StatusNotFound, "商品不存在")
		}
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return data, nil
}

// GetCategoryList 分类列表
// @Summary      分类列表
// @Description  分页查询商品分类
// @Tags         分类
// @Produce      json
// @Param        page       query  int  false  "页码"       default(1)
// @Param        page_size  query  int  false  "每页数量"   default(10)
// @Success      200  {object}  apidoc.Response{data=dto.CategoryListResp}  "查询成功"
// @Router       /api/v1/product_category/list [get]
func (h *CatalogHandler) GetCategoryList(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	data, err := h.logic.GetCategoryList(ctx, pageReq)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return data, nil
}

// GetCategoryDetail 分类详情
// @Summary      分类详情
// @Description  根据分类 ID 查询详情
// @Tags         分类
// @Produce      json
// @Param        id  query  int  true  "分类 ID"
// @Success      200  {object}  apidoc.Response{data=dto.CategoryDetailResp}  "查询成功"
// @Router       /api/v1/product_category/detail [get]
func (h *CatalogHandler) GetCategoryDetail(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.QueryGet("id"), 10, 64)
	if id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	data, err := h.logic.GetCategoryDetail(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerr.New(http.StatusNotFound, "分类不存在")
		}
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return data, nil
}
