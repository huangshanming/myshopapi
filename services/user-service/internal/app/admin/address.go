package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"
)

func (h *AddressHandler) AdminList(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "用户ID无效")
	}
	list, err := h.logic.List(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}
