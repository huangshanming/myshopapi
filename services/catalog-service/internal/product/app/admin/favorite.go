package admin

import (
	"mymall/pkg/httpserver"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *FavoriteHandler) AdminUserList(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "用户ID无效"))
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 50
	}
	list, total, err := h.logic.List(r.Context(), userID, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}
