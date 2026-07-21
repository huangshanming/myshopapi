package user

import (
	"encoding/json"
	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *FavoriteHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	var req struct {
		ProductID uint64 `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProductID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.Add(r.Context(), userID, req.ProductID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *FavoriteHandler) Remove(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	productID, err := strconv.ParseUint(httpserver.PathParam(r, "product_id"), 10, 64)
	if err != nil || productID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	if err := h.logic.Remove(r.Context(), userID, productID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *FavoriteHandler) RemoveBatch(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	var req struct {
		ProductIDs []uint64 `json:"product_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.ProductIDs) == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.RemoveBatch(r.Context(), userID, req.ProductIDs); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *FavoriteHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.List(r.Context(), userID, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *FavoriteHandler) Status(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	productID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || productID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	okFav, err := h.logic.IsFavorited(r.Context(), userID, productID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]bool{"favorited": okFav})
}

func (h *FavoriteHandler) Count(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || productID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	n, err := h.logic.FavoriteCount(productID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "商品不存在"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]int64{"count": n})
}
