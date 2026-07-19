package handler

import (
	"net/http"

	"mymall/pkg/response"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

type RegionHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.RegionLogic
}

func NewRegionHandler(svcCtx *svc.ServiceContext) *RegionHandler {
	return &RegionHandler{
		svcCtx: svcCtx,
		logic:  logic.NewRegionLogic(svcCtx),
	}
}

func (h *RegionHandler) List(w http.ResponseWriter, r *http.Request) {
	parent := r.URL.Query().Get("parent_code")
	list, err := h.logic.ListChildren(parent)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "ok")
}

func (h *RegionHandler) Tree(w http.ResponseWriter, r *http.Request) {
	tree, err := h.logic.Tree()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, tree, "ok")
}
