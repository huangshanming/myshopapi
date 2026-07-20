package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"mymall/pkg/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type addressDeps struct {
	svcCtx *svc.ServiceContext
	logic  *logic.AddressLogic
}

func newAddressDeps(svcCtx *svc.ServiceContext) addressDeps {
	return addressDeps{
		svcCtx: svcCtx,
		logic:  logic.NewAddressLogic(context.Background(), svcCtx),
	}
}

type AddressUserHandler struct{ addressDeps }
type AddressAdminHandler struct{ addressDeps }
type AddressInternalHandler struct{ addressDeps }

func NewAddressUserHandler(svcCtx *svc.ServiceContext) *AddressUserHandler {
	return &AddressUserHandler{addressDeps: newAddressDeps(svcCtx)}
}
func NewAddressAdminHandler(svcCtx *svc.ServiceContext) *AddressAdminHandler {
	return &AddressAdminHandler{addressDeps: newAddressDeps(svcCtx)}
}
func NewAddressInternalHandler(svcCtx *svc.ServiceContext) *AddressInternalHandler {
	return &AddressInternalHandler{addressDeps: newAddressDeps(svcCtx)}
}

func (h *AddressUserHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	list, err := h.logic.List(userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *AddressUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	var req types.AddressReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	a, err := h.logic.Create(userID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, a)
}

func (h *AddressUserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "地址ID无效"))
		return
	}
	var req types.AddressReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.Update(userID, id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AddressUserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "地址ID无效"))
		return
	}
	if err := h.logic.Delete(userID, id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AddressUserHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "地址ID无效"))
		return
	}
	if err := h.logic.SetDefault(userID, id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AddressAdminHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "用户ID无效"))
		return
	}
	list, err := h.logic.List(userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *AddressInternalHandler) InternalGet(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	id, _ := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if userID == 0 || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数无效"))
		return
	}
	a, err := h.logic.Get(userID, id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, a)
}
