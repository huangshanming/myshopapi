package admin

import (
	"encoding/json"
	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *CatalogHandler) AdminListProducts(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	status := r.URL.Query().Get("status")
	data, err := h.logic.GetProductListFiltered(r.Context(), pageReq, shopID, status, 0, "")
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *CatalogHandler) AdminForceOffSale(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	if err := h.logic.ForceOffSale(r.Context(), id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CatalogHandler) AdminListCategories(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListAllCategories(r.Context())
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *CatalogHandler) AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req types.CategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if req.Name == "" {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	cat, err := h.logic.CreateCategory(r.Context(), req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, cat)
}

func (h *CatalogHandler) AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
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
	if err := h.logic.UpdateCategory(r.Context(), id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CatalogHandler) AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "分类ID无效"))
		return
	}
	if err := h.logic.DeleteCategory(r.Context(), id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
