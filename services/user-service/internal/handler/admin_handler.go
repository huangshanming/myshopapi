package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type AdminHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.RBACLogic
}

func NewAdminHandler(svcCtx *svc.ServiceContext) *AdminHandler {
	return &AdminHandler{
		svcCtx: svcCtx,
		logic:  logic.NewRBACLogic(svcCtx),
	}
}

// PermChecker for middleware
func (h *AdminHandler) IsSuperAdmin(userID uint64) bool { return h.logic.IsSuperAdmin(userID) }
func (h *AdminHandler) HasPerm(userID uint64, code string) bool {
	return h.logic.HasPerm(userID, code)
}

func (h *AdminHandler) AuthMe(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r.Context())
	data, err := h.logic.AuthMe(uid)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *AdminHandler) MenuTree(w http.ResponseWriter, r *http.Request) {
	tree, err := h.logic.MenuTreeAll()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, tree, "ok")
}

func (h *AdminHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var req types.MenuReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	m, err := h.logic.CreateMenu(req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, m, "创建成功")
}

func (h *AdminHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.MenuReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.UpdateMenu(id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *AdminHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.DeleteMenu(id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "删除成功")
}

func (h *AdminHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListRoles()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "ok")
}

func (h *AdminHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req types.RoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	role, err := h.logic.CreateRole(req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, role, "创建成功")
}

func (h *AdminHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.RoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.UpdateRole(id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *AdminHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.DeleteRole(id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "删除成功")
}

func (h *AdminHandler) GetRoleMenus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	ids, err := h.logic.GetRoleMenus(id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, ids, "ok")
}

func (h *AdminHandler) AssignRoleMenus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.RoleMenusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.AssignRoleMenus(id, req.MenuIDs); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "分配成功")
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	mobile := r.URL.Query().Get("mobile")
	list, total, err := h.logic.ListUsers(page, pageSize, mobile)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "ok")
}

func (h *AdminHandler) SetUserStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.UserStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.SetUserStatus(id, req.Status); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	user, err := h.logic.GetUser(id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response.Success(w, user, "ok")
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.UserUpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.UpdateUser(id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *AdminHandler) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.UserResetPwdReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.ResetUserPassword(id, req.Password); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "重置成功")
}

func (h *AdminHandler) GenerateUserToken(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	data, err := h.logic.GenerateUserToken(id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, data, "ok")
}

func (h *AdminHandler) ListAdmins(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.ListAdmins(page, pageSize)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "ok")
}

func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req types.AdminCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	user, err := h.logic.CreateAdmin(req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, user, "创建成功")
}

func (h *AdminHandler) GetAdminRoles(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	ids, err := h.logic.AdminRoleIDs(id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, ids, "ok")
}

func (h *AdminHandler) AssignAdminRoles(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.AdminRolesReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.AssignAdminRoles(id, req.RoleIDs); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "分配成功")
}

func (h *AdminHandler) ResetAdminPassword(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.AdminResetPwdReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.ResetAdminPassword(id, req.Password); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "重置成功")
}

func (h *AdminHandler) ListConfigs(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListConfigs()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "ok")
}

func (h *AdminHandler) SaveConfigs(w http.ResponseWriter, r *http.Request) {
	var req types.ConfigBatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.SaveConfigs(req.Items); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "保存成功")
}
