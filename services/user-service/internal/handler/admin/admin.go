package admin

import (
	"encoding/json"
	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/types"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *AdminHandler) AdminSendNotification(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())
	var req logic.AdminSendReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	userLogic := logic.NewUserLogic(r.Context(), h.svcCtx)
	batch, err := userLogic.AdminSend(adminID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, batch)
}

func (h *AdminHandler) AdminListNotificationSends(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	userLogic := logic.NewUserLogic(r.Context(), h.svcCtx)
	list, total, err := userLogic.ListSendBatches(page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *AdminHandler) AdminListNotificationRecipients(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	page, pageSize := middleware.ParsePage(r)
	userLogic := logic.NewUserLogic(r.Context(), h.svcCtx)
	batch, err := userLogic.GetSendBatch(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "发送记录不存在"))
		return
	}
	list, total, err := userLogic.ListBatchRecipients(id, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"batch": batch,
		"list":  list,
		"total": total,
	})
}

func (h *AdminHandler) IsSuperAdmin(userID uint64) bool { return h.logic.IsSuperAdmin(userID) }

func (h *AdminHandler) HasPerm(userID uint64, code string) bool {
	return h.logic.HasPerm(userID, code)
}

func (h *AdminHandler) AuthMe(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r.Context())
	data, err := h.logic.AuthMe(uid)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *AdminHandler) MenuTree(w http.ResponseWriter, r *http.Request) {
	tree, err := h.logic.MenuTreeAll()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, tree)
}

func (h *AdminHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var req types.MenuReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	m, err := h.logic.CreateMenu(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, m)
}

func (h *AdminHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.MenuReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateMenu(id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.DeleteMenu(id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListRoles()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *AdminHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req types.RoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	role, err := h.logic.CreateRole(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, role)
}

func (h *AdminHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.RoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateRole(id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.DeleteRole(id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) GetRoleMenus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	ids, err := h.logic.GetRoleMenus(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, ids)
}

func (h *AdminHandler) AssignRoleMenus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.RoleMenusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.AssignRoleMenus(id, req.MenuIDs); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	mobile := r.URL.Query().Get("mobile")
	list, total, err := h.logic.ListUsers(page, pageSize, mobile)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *AdminHandler) SetUserStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.UserStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.SetUserStatus(id, req.Status); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	user, err := h.logic.GetUser(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, user)
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.UserUpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateUser(id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.UserResetPwdReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.ResetUserPassword(id, req.Password); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) GenerateUserToken(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	data, err := h.logic.GenerateUserToken(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *AdminHandler) ListAdmins(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.ListAdmins(page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req types.AdminCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	user, err := h.logic.CreateAdmin(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, user)
}

func (h *AdminHandler) GetAdminRoles(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	ids, err := h.logic.AdminRoleIDs(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, ids)
}

func (h *AdminHandler) AssignAdminRoles(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.AdminRolesReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.AssignAdminRoles(id, req.RoleIDs); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) ResetAdminPassword(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.AdminResetPwdReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.ResetAdminPassword(id, req.Password); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *AdminHandler) ListConfigs(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListConfigs()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *AdminHandler) SaveConfigs(w http.ResponseWriter, r *http.Request) {
	var req types.ConfigBatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.SaveConfigs(req.Items); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
