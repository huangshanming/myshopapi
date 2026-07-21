package public

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"net/http"
)

func (h *RegionHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	parent := in.QueryGet("parent_code")
	list, err := h.logic.ListChildren(ctx, parent)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *RegionHandler) Tree(ctx context.Context, in appinput.CallInput) (any, error) {
	tree, err := h.logic.Tree(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return tree, nil
}
