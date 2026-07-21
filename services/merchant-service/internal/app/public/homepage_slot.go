package public

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	"mymall/pkg/xerr"
)

func (h *HomepageSlotHandler) PublicHomeSlots(ctx context.Context, in appinput.CallInput) (any, error) {
	slotType := in.QueryGet("slot_type")
	list, err := h.logic.HomeSlots(slotType)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return list, nil
}
