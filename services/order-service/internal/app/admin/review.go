package admin

import (
	"mymall/pkg/httpserver"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *ReviewHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	level := r.URL.Query().Get("rating_level")
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	list, total, err := h.logic.AdminList(r.Context(), shopID, level, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *ReviewHandler) AdminDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "评价ID无效"))
		return
	}
	if err := h.logic.SoftDelete(r.Context(), id, 0); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
