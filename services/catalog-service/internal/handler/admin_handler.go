package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/response"
	"mymall/services/catalog-service/internal/model"

	"github.com/gin-gonic/gin"
)

type merchantProductReq struct {
	Name       string  `json:"name" binding:"required"`
	SalePrice  float64 `json:"sale_price" binding:"required"`
	Stock      int     `json:"stock"`
	CategoryID uint64  `json:"category_id" binding:"required"`
	Subtitle   string  `json:"subtitle"`
	MainImage  string  `json:"main_image"`
	Status     string  `json:"status"`
	PetType    string  `json:"pet_type"`
}

func (h *CatalogHandler) MerchantListProducts(c *gin.Context) {
	shopID := middleware.GetShopID(c)
	if shopID == 0 {
		response.Error(c, "缺少 shop_id", http.StatusForbidden)
		return
	}
	pageReq := &pagination.PageReq{}
	_ = c.ShouldBindQuery(pageReq)
	p, ps, _ := pagination.Normalize(pageReq)
	pageReq.Page, pageReq.PageSize = p, ps
	status := c.Query("status")
	data, err := h.svc.GetProductListFiltered(pageReq, shopID, status)
	if err != nil {
		response.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(c, data, "查询成功")
}

func (h *CatalogHandler) MerchantCreateProduct(c *gin.Context) {
	shopID := middleware.GetShopID(c)
	if shopID == 0 {
		response.Error(c, "缺少 shop_id", http.StatusForbidden)
		return
	}
	var req merchantProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	status := req.Status
	if status == "" {
		status = "on_sale"
	}
	pet := req.PetType
	if pet == "" {
		pet = "both"
	}
	p := &model.Product{
		ShopID:     shopID,
		ProductNo:  fmt.Sprintf("P%d", time.Now().UnixNano()%1e12),
		Name:       req.Name,
		SalePrice:  req.SalePrice,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
		Subtitle:   req.Subtitle,
		MainImage:  req.MainImage,
		Status:     status,
		PetType:    pet,
		Discount:   100,
	}
	if err := h.svc.CreateProduct(p); err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, p, "创建成功")
}

func (h *CatalogHandler) MerchantUpdateProduct(c *gin.Context) {
	shopID := middleware.GetShopID(c)
	if shopID == 0 {
		response.Error(c, "缺少 shop_id", http.StatusForbidden)
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "商品ID无效", http.StatusBadRequest)
		return
	}
	var req merchantProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	updates := map[string]interface{}{
		"name":        req.Name,
		"sale_price":  req.SalePrice,
		"stock":       req.Stock,
		"category_id": req.CategoryID,
		"subtitle":    req.Subtitle,
		"main_image":  req.MainImage,
	}
	if req.Status != "" {
		updates["status"] = req.Status
		if req.Status == "on_sale" {
			updates["publish_time"] = time.Now()
		}
	}
	if req.PetType != "" {
		updates["pet_type"] = req.PetType
	}
	if err := h.svc.UpdateProductByShop(id, shopID, updates); err != nil {
		response.Error(c, err.Error(), http.StatusForbidden)
		return
	}
	response.Success(c, nil, "更新成功")
}

func (h *CatalogHandler) MerchantSetStatus(c *gin.Context) {
	shopID := middleware.GetShopID(c)
	if shopID == 0 {
		response.Error(c, "缺少 shop_id", http.StatusForbidden)
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "商品ID无效", http.StatusBadRequest)
		return
	}
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.svc.UpdateProductByShop(id, shopID, map[string]interface{}{"status": body.Status}); err != nil {
		response.Error(c, err.Error(), http.StatusForbidden)
		return
	}
	response.Success(c, nil, "更新成功")
}

func (h *CatalogHandler) AdminListProducts(c *gin.Context) {
	pageReq := &pagination.PageReq{}
	_ = c.ShouldBindQuery(pageReq)
	p, ps, _ := pagination.Normalize(pageReq)
	pageReq.Page, pageReq.PageSize = p, ps
	shopID, _ := strconv.ParseUint(c.Query("shop_id"), 10, 64)
	status := c.Query("status")
	data, err := h.svc.GetProductListFiltered(pageReq, shopID, status)
	if err != nil {
		response.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(c, data, "查询成功")
}

func (h *CatalogHandler) AdminForceOffSale(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "商品ID无效", http.StatusBadRequest)
		return
	}
	if err := h.svc.ForceOffSale(id); err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, nil, "已强制下架")
}

type categoryReq struct {
	ParentId    uint64 `json:"parent_id"`
	Name        string `json:"name" binding:"required"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Level       int    `json:"level"`
	IsShow      *bool  `json:"is_show"`
}

func (h *CatalogHandler) AdminCreateCategory(c *gin.Context) {
	var req categoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	level := req.Level
	if level == 0 {
		level = 1
	}
	show := true
	if req.IsShow != nil {
		show = *req.IsShow
	}
	cat := &model.ProductCategory{
		ParentId:    req.ParentId,
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Level:       level,
		IsShow:      show,
	}
	if err := h.svc.CreateCategory(cat); err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, cat, "创建成功")
}

func (h *CatalogHandler) AdminUpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "分类ID无效", http.StatusBadRequest)
		return
	}
	var req categoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	updates := map[string]interface{}{
		"parent_id":   req.ParentId,
		"name":        req.Name,
		"icon":        req.Icon,
		"description": req.Description,
		"sort_order":  req.SortOrder,
	}
	if req.Level > 0 {
		updates["level"] = req.Level
	}
	if req.IsShow != nil {
		updates["is_show"] = *req.IsShow
	}
	if err := h.svc.UpdateCategory(id, updates); err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, nil, "更新成功")
}

func (h *CatalogHandler) AdminDeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "分类ID无效", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteCategory(id); err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, nil, "删除成功")
}

// Ensure jwt roles imported for route wiring docs
var _ = jwt.RolePlatformAdmin
