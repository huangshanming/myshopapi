package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/apidoc/dto"
	"mymall/pkg/pagination"
	"mymall/pkg/response"
	"mymall/services/catalog-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CatalogHandler struct {
	svc *service.CatalogService
}

// swag 类型引用
var (
	_ dto.ProductListResp
	_ dto.ProductDetail
	_ dto.CategoryListResp
	_ dto.CategoryDetailResp
)

func NewCatalogHandler(svc *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{svc: svc}
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
func (h *CatalogHandler) GetProductList(c *gin.Context) {
	pageReq := c.MustGet("pageReq").(*pagination.PageReq)
	data, err := h.svc.GetProductList(pageReq)
	if err != nil {
		response.Error(c, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(c, data, "查询成功")
}

// GetProductDetail 商品详情
// @Summary      商品详情
// @Description  根据商品 ID 查询详情
// @Tags         商品
// @Produce      json
// @Param        id  query  int  true  "商品 ID"
// @Success      200  {object}  apidoc.Response{data=dto.ProductDetail}  "查询成功"
// @Router       /api/v1/products/detail [get]
func (h *CatalogHandler) GetProductDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	data, err := h.svc.GetProductDetail(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, "商品不存在", http.StatusNotFound)
			return
		}
		response.Error(c, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(c, data, "查询成功")
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
func (h *CatalogHandler) GetCategoryList(c *gin.Context) {
	pageReq := c.MustGet("pageReq").(*pagination.PageReq)
	data, err := h.svc.GetCategoryList(pageReq)
	if err != nil {
		response.Error(c, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(c, data, "查询成功")
}

// GetCategoryDetail 分类详情
// @Summary      分类详情
// @Description  根据分类 ID 查询详情
// @Tags         分类
// @Produce      json
// @Param        id  query  int  true  "分类 ID"
// @Success      200  {object}  apidoc.Response{data=dto.CategoryDetailResp}  "查询成功"
// @Router       /api/v1/product_category/detail [get]
func (h *CatalogHandler) GetCategoryDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	data, err := h.svc.GetCategoryDetail(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, "分类不存在", http.StatusNotFound)
			return
		}
		response.Error(c, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(c, data, "查询成功")
}
