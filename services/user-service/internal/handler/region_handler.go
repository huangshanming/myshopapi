package handler

import (
	"context"
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/pkg/xerr")

type RegionHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.RegionLogic
}

func NewRegionHandler(svcCtx *svc.ServiceContext) *RegionHandler {
	return &RegionHandler{
		svcCtx: svcCtx,
		logic:  logic.NewRegionLogic(context.Background(), svcCtx),
	}
}

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
