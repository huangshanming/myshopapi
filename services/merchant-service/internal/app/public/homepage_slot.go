package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/xerr"
)

func (h *HomepageSlotHandler) PublicHomeSlots(w http.ResponseWriter, r *http.Request) {
	slotType := r.URL.Query().Get("slot_type")
	list, err := h.logic.HomeSlots(slotType)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}
