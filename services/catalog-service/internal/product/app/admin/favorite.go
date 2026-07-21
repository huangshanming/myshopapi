package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"
)

func (h *FavoriteHandler) AdminUserList(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || userID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "用户ID无效")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	if pageSize <= 0 {
		pageSize = 50
	}
	list, total, err := h.logic.List(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}
