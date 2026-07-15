package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/response"
	"mymall/services/catalog-service/internal/types"
)

func (h *CatalogHandler) MerchantListProducts(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少 shop_id", http.StatusForbidden)
		return
	}
	page, pageSize := middleware.ParsePage(r)
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	status := r.URL.Query().Get("status")
	data, err := h.logic.GetProductListFiltered(pageReq, shopID, status)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "查询成功")
}

func (h *CatalogHandler) MerchantCreateProduct(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少 shop_id", http.StatusForbidden)
		return
	}
	var req types.MerchantProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.SalePrice == 0 || req.CategoryID == 0 {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	p, err := h.logic.CreateProduct(shopID, req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, p, "创建成功")
}

func (h *CatalogHandler) MerchantUpdateProduct(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少 shop_id", http.StatusForbidden)
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "商品ID无效", http.StatusBadRequest)
		return
	}
	var req types.MerchantProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.SalePrice == 0 || req.CategoryID == 0 {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.UpdateProductByShop(id, shopID, req); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *CatalogHandler) MerchantSetStatus(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少 shop_id", http.StatusForbidden)
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "商品ID无效", http.StatusBadRequest)
		return
	}
	var body types.SetStatusReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.SetProductStatus(id, shopID, body.Status); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *CatalogHandler) AdminListProducts(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	status := r.URL.Query().Get("status")
	data, err := h.logic.GetProductListFiltered(pageReq, shopID, status)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "查询成功")
}

func (h *CatalogHandler) AdminForceOffSale(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "商品ID无效", http.StatusBadRequest)
		return
	}
	if err := h.logic.ForceOffSale(id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已强制下架")
}

func (h *CatalogHandler) AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req types.CategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	cat, err := h.logic.CreateCategory(req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, cat, "创建成功")
}

func (h *CatalogHandler) AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "分类ID无效", http.StatusBadRequest)
		return
	}
	var req types.CategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.UpdateCategory(id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *CatalogHandler) AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "分类ID无效", http.StatusBadRequest)
		return
	}
	if err := h.logic.DeleteCategory(id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "删除成功")
}
