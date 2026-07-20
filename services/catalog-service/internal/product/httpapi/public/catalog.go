package public

import (
	"mymall/pkg/apidoc/dto"
	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
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
func (h *CatalogHandler) GetProductList(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	categoryID, _ := strconv.ParseUint(r.URL.Query().Get("category_id"), 10, 64)
	orderBy := r.URL.Query().Get("order_by")
	data, err := h.logic.GetProductListFiltered(pageReq, shopID, "on_sale", categoryID, orderBy)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, "查询失败"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

// GetSalesRank 今日必买销量榜
func (h *CatalogHandler) GetSalesRank(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	data, err := h.logic.GetSalesRank(page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, "查询失败"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

// GetProductDetail 商品详情
// @Summary      商品详情
// @Description  根据商品 ID 查询详情
// @Tags         商品
// @Produce      json
// @Param        id  query  int  true  "商品 ID"
// @Success      200  {object}  apidoc.Response{data=dto.ProductDetail}  "查询成功"
// @Router       /api/v1/products/detail [get]
func (h *CatalogHandler) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	data, err := h.logic.GetProductDetail(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "商品不存在"))
			return
		}
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, "查询失败"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
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
func (h *CatalogHandler) GetCategoryList(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	data, err := h.logic.GetCategoryList(pageReq)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, "查询失败"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

// GetCategoryDetail 分类详情
// @Summary      分类详情
// @Description  根据分类 ID 查询详情
// @Tags         分类
// @Produce      json
// @Param        id  query  int  true  "分类 ID"
// @Success      200  {object}  apidoc.Response{data=dto.CategoryDetailResp}  "查询成功"
// @Router       /api/v1/product_category/detail [get]
func (h *CatalogHandler) GetCategoryDetail(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	data, err := h.logic.GetCategoryDetail(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "分类不存在"))
			return
		}
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, "查询失败"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}
