package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/pagination"
	"mymall/pkg/response"
	"mymall/services/catalog-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CatalogHandler struct {
	svc *service.CatalogService
}

func NewCatalogHandler(svc *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{svc: svc}
}

func (h *CatalogHandler) GetProductList(c *gin.Context) {
	pageReq := c.MustGet("pageReq").(*pagination.PageReq)
	data, err := h.svc.GetProductList(pageReq)
	if err != nil {
		response.Error(c, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(c, data, "查询成功")
}

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

func (h *CatalogHandler) GetCategoryList(c *gin.Context) {
	pageReq := c.MustGet("pageReq").(*pagination.PageReq)
	data, err := h.svc.GetCategoryList(pageReq)
	if err != nil {
		response.Error(c, "查询失败", http.StatusInternalServerError)
		return
	}
	response.Success(c, data, "查询成功")
}

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
