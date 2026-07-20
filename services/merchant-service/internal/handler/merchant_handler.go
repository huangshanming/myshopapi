package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *ShopMerchantHandler) Apply(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	var in types.ApplyReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	app, err := h.logic.Apply(userID, in)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, app)
}

func (h *ShopMerchantHandler) MyShops(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	shops, err := h.logic.MyShops(userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, shops)
}

func (h *ShopMerchantHandler) UpdateMyShop(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	shopID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	var req types.UpdateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateMyShop(shopID, userID, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ShopAdminHandler) AdminListApplications(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListApplications(r.URL.Query().Get("status"), p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *ShopAdminHandler) AdminApprove(w http.ResponseWriter, r *http.Request) {
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
	shop, err := h.logic.Approve(appID, adminID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, shop)
}

func (h *ShopAdminHandler) AdminReject(w http.ResponseWriter, r *http.Request) {
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
	if err := h.logic.Reject(appID, adminID, req.Reason); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ShopAdminHandler) AdminListShops(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListShops(r.URL.Query().Get("status"), r.URL.Query().Get("name"), p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

// PublicListShops C 端公开商户列表（无需登录）
func (h *ShopPublicHandler) PublicListShops(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListPublicShops(p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *ShopPublicHandler) PublicGetShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	shop, err := h.logic.GetPublicShop(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, shop)
}

func (h *ShopAdminHandler) AdminCreateShop(w http.ResponseWriter, r *http.Request) {
	var req types.AdminCreateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	shop, err := h.logic.CreateShop(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, shop)
}

func (h *ShopAdminHandler) AdminUpdateShop(w http.ResponseWriter, r *http.Request) {
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
	if err := h.logic.AdminUpdateShop(id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ShopAdminHandler) AdminResetOwnerPassword(w http.ResponseWriter, r *http.Request) {
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
	if err := h.logic.ResetOwnerPassword(id, req.Password); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ShopAdminHandler) AdminGetShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	shop, err := h.logic.GetShop(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "店铺不存在"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, shop)
}

func (h *ShopAdminHandler) AdminDisableShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	var req types.RejectReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.DisableShop(id, req.Reason); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ShopAdminHandler) AdminEnableShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	if err := h.logic.EnableShop(id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
