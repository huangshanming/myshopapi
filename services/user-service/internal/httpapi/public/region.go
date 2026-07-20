package public

import (
	"mymall/pkg/xerr"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *RegionHandler) List(w http.ResponseWriter, r *http.Request) {
	parent := r.URL.Query().Get("parent_code")
	list, err := h.logic.ListChildren(parent)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *RegionHandler) Tree(w http.ResponseWriter, r *http.Request) {
	tree, err := h.logic.Tree()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, tree)
}
