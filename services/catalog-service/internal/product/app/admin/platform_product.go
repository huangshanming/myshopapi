package admin

import (
	"encoding/json"
	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *PlatformProductHandler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	catID, _ := strconv.ParseUint(r.URL.Query().Get("category_id"), 10, 64)
	f := repository.ProductListFilter{
		ShopID: shopID, Name: r.URL.Query().Get("name"), ProductNo: r.URL.Query().Get("product_no"),
		CategoryID: catID, Status: r.URL.Query().Get("status"), ProductType: r.URL.Query().Get("product_type"),
		Page: page, PageSize: pageSize, OrderBy: r.URL.Query().Get("order_by"),
		PlatformScope: true,
	}
	if s := r.URL.Query().Get("created_from"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if s := r.URL.Query().Get("created_to"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	if s := r.URL.Query().Get("publish_from"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.PublishFrom = &t
		}
	}
	if s := r.URL.Query().Get("publish_to"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.PublishTo = &end
		}
	}
	data, err := h.logic.List(r.Context(), f)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *PlatformProductHandler) OffSale(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r.Context())
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.PlatformProductRemarkReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.ForceOffSale(r.Context(), id, uid, req.Remark); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *PlatformProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r.Context())
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.PlatformProductRemarkReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.SoftDelete(r.Context(), id, uid, req.Remark); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
