package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/apidoc/dto"
	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/response"
	"mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/svc"

	"gorm.io/gorm"
)

type CatalogHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.CatalogLogic
}

// swag 类型引用
var (
	_ dto.ProductListResp
	_ dto.ProductDetail
	_ dto.CategoryListResp
	_ dto.CategoryDetailResp
)

func NewCatalogHandler(svcCtx *svc.ServiceContext) *CatalogHandler {
	return &CatalogHandler{
		svcCtx: svcCtx,
		logic:  logic.NewCatalogLogic(svcCtx),
	}
}

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
	data, err := h.logic.GetProductList(pageReq)
	if err != nil {
		response.Error(w, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "查询成功")
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
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	data, err := h.logic.GetProductDetail(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(w, "商品不存在", http.StatusNotFound)
			return
		}
		response.Error(w, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "查询成功")
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
		response.Error(w, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "查询成功")
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
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	data, err := h.logic.GetCategoryDetail(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(w, "分类不存在", http.StatusNotFound)
			return
		}
		response.Error(w, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "查询成功")
}
