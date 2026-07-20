package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/xerr"
)

func (h *HomepageThemeHandler) PublicThemeTiles(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListThemeTiles()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list})
}
