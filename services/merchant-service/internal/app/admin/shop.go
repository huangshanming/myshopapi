package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *ShopHandler) AdminListApplications(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListApplications(r.Context(), r.URL.Query().Get("status"), p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *ShopHandler) AdminApprove(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	appID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "申请ID无效"))
		return
	}
	shop, err := h.logic.Approve(r.Context(), appID, adminID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, shop)
}

func (h *ShopHandler) AdminReject(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	appID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "申请ID无效"))
		return
	}
	var req types.RejectReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.Reject(r.Context(), appID, adminID, req.Reason); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ShopHandler) AdminListShops(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListShops(r.Context(), r.URL.Query().Get("status"), r.URL.Query().Get("name"), p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *ShopHandler) AdminCreateShop(w http.ResponseWriter, r *http.Request) {
	var req types.AdminCreateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	shop, err := h.logic.CreateShop(r.Context(), req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, shop)
}

func (h *ShopHandler) AdminUpdateShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	var req types.AdminUpdateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.AdminUpdateShop(r.Context(), id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ShopHandler) AdminResetOwnerPassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	var req types.OwnerPasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.ResetOwnerPassword(r.Context(), id, req.Password); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ShopHandler) AdminGetShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	shop, err := h.logic.GetShop(r.Context(), id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "店铺不存在"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, shop)
}

func (h *ShopHandler) AdminDisableShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	var req types.RejectReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.DisableShop(r.Context(), id, req.Reason); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ShopHandler) AdminEnableShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	if err := h.logic.EnableShop(r.Context(), id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
