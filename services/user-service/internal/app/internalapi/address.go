package internalapi

import (
	"mymall/pkg/xerr"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *AddressHandler) InternalGet(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	id, _ := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if userID == 0 || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数无效"))
		return
	}
	a, err := h.logic.Get(r.Context(), userID, id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, a)
}
