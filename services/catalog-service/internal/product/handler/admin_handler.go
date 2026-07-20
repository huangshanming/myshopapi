package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/services/catalog-service/internal/product/types"

	"mymall/pkg/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *CatalogMerchantHandler) MerchantListProducts(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少 shop_id"))
		return
	}
	page, pageSize := middleware.ParsePage(r)
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	status := r.URL.Query().Get("status")
	data, err := h.logic.GetProductListFiltered(pageReq, shopID, status, 0, "")
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *CatalogMerchantHandler) MerchantCreateProduct(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少 shop_id"))
		return
	}
	var req types.MerchantProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if req.Name == "" || req.SalePrice == 0 || req.CategoryID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	p, err := h.logic.CreateProduct(shopID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *CatalogMerchantHandler) MerchantUpdateProduct(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少 shop_id"))
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	var req types.MerchantProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if req.Name == "" || req.SalePrice == 0 || req.CategoryID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateProductByShop(id, shopID, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CatalogMerchantHandler) MerchantSetStatus(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少 shop_id"))
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	var body types.SetStatusReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.SetProductStatus(id, shopID, body.Status); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CatalogAdminHandler) AdminListProducts(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	status := r.URL.Query().Get("status")
	data, err := h.logic.GetProductListFiltered(pageReq, shopID, status, 0, "")
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *CatalogAdminHandler) AdminForceOffSale(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	if err := h.logic.ForceOffSale(id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CatalogAdminHandler) AdminListCategories(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListAllCategories()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *CatalogAdminHandler) AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req types.CategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if req.Name == "" {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	cat, err := h.logic.CreateCategory(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, cat)
}

func (h *CatalogAdminHandler) AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "分类ID无效"))
		return
	}
	var req types.CategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if req.Name == "" {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateCategory(id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CatalogAdminHandler) AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "分类ID无效"))
		return
	}
	if err := h.logic.DeleteCategory(id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
