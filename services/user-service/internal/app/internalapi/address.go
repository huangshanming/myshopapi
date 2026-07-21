package internalapi

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"
)

func (h *AddressHandler) InternalGet(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, _ := strconv.ParseUint(in.QueryGet("user_id"), 10, 64)
	id, _ := strconv.ParseUint(in.QueryGet("id"), 10, 64)
	if userID == 0 || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数无效")
	}
	a, err := h.logic.Get(ctx, userID, id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return a, nil
}
