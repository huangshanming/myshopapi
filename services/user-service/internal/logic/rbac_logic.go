package logic

import (
	"errors"
	"strings"

	"mymall/pkg/jwt"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type RBACLogic struct {
	svcCtx *svc.ServiceContext
}

func NewRBACLogic(svcCtx *svc.ServiceContext) *RBACLogic {
	return &RBACLogic{svcCtx: svcCtx}
}

func (l *RBACLogic) IsSuperAdmin(userID uint64) bool {
	return l.svcCtx.RBAC.IsSuperAdmin(userID)
}

func (l *RBACLogic) HasPerm(userID uint64, code string) bool {
	perms, err := l.svcCtx.RBAC.ListUserPerms(userID)
	if err != nil {
		return false
	}
	for _, p := range perms {
		if p == code {
			return true
		}
	}
	return false
}

func BuildMenuTree(menus []model.SysMenu, onlyVisible bool) []*types.MenuTreeNode {
	nodes := make(map[uint64]*types.MenuTreeNode, len(menus))
	var roots []*types.MenuTreeNode
	for _, m := range menus {
		if onlyVisible && m.Visible != 1 {
			continue
		}
		n := &types.MenuTreeNode{
			ID: m.ID, ParentID: m.ParentID, Name: m.Name, Type: m.Type,
			Path: m.Path, Component: m.Component, Icon: m.Icon, Perms: m.Perms,
			Sort: m.Sort, Visible: m.Visible, Status: m.Status,
			Children: []*types.MenuTreeNode{},
		}
		nodes[m.ID] = n
	}
	for _, m := range menus {
		n, ok := nodes[m.ID]
		if !ok {
			continue
		}
		if m.ParentID == 0 {
			roots = append(roots, n)
			continue
		}
		if p, ok := nodes[m.ParentID]; ok {
			p.Children = append(p.Children, n)
		} else {
			roots = append(roots, n)
		}
	}
	return roots
}

func (l *RBACLogic) AuthMe(userID uint64) (*types.AuthMeResp, error) {
	roles, err := l.svcCtx.RBAC.ListRoleCodes(userID)
	if err != nil {
		return nil, err
	}
	perms, err := l.svcCtx.RBAC.ListUserPerms(userID)
	if err != nil {
		return nil, err
	}
	menus, err := l.svcCtx.RBAC.ListUserMenus(userID)
	if err != nil {
		return nil, err
	}
	return &types.AuthMeResp{
		Roles:    roles,
		Perms:    perms,
		MenuTree: BuildMenuTree(menus, true),
	}, nil
}

func (l *RBACLogic) MenuTreeAll() ([]*types.MenuTreeNode, error) {
	menus, err := l.svcCtx.RBAC.ListAllMenus()
	if err != nil {
		return nil, err
	}
	return BuildMenuTree(menus, false), nil
}

func (l *RBACLogic) CreateMenu(req types.MenuReq) (*model.SysMenu, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("名称不能为空")
	}
	typ := req.Type
	if typ == "" {
		typ = model.MenuTypeMenu
	}
	m := &model.SysMenu{
		ParentID: req.ParentID, Name: req.Name, Type: typ,
		Path: req.Path, Component: req.Component, Icon: req.Icon, Perms: req.Perms,
		Sort: req.Sort, Visible: req.Visible, Status: req.Status,
	}
	if m.Visible == 0 && req.Visible == 0 {
		m.Visible = 1
	}
	if m.Status == 0 {
		m.Status = 1
	}
	if err := l.svcCtx.RBAC.CreateMenu(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (l *RBACLogic) UpdateMenu(id uint64, req types.MenuReq) error {
	m, err := l.svcCtx.RBAC.GetMenu(id)
	if err != nil {
		return errors.New("菜单不存在")
	}
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.Type != "" {
		m.Type = req.Type
	}
	m.ParentID = req.ParentID
	m.Path = req.Path
	m.Component = req.Component
	m.Icon = req.Icon
	m.Perms = req.Perms
	m.Sort = req.Sort
	m.Visible = req.Visible
	m.Status = req.Status
	if m.Status == 0 && req.Status == 0 {
		// allow disabled
	}
	return l.svcCtx.RBAC.UpdateMenu(m)
}

func (l *RBACLogic) DeleteMenu(id uint64) error {
	if err := l.svcCtx.RBAC.DeleteMenu(id); err != nil {
		return errors.New("请先删除子菜单")
	}
	return nil
}

func (l *RBACLogic) ListRoles() ([]model.SysRole, error) {
	return l.svcCtx.RBAC.ListRoles()
}

func (l *RBACLogic) CreateRole(req types.RoleReq) (*model.SysRole, error) {
	if req.Code == "" || req.Name == "" {
		return nil, errors.New("编码和名称不能为空")
	}
	role := &model.SysRole{Code: req.Code, Name: req.Name, Status: req.Status, Remark: req.Remark}
	if role.Status == 0 {
		role.Status = 1
	}
	if err := l.svcCtx.RBAC.CreateRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (l *RBACLogic) UpdateRole(id uint64, req types.RoleReq) error {
	role, err := l.svcCtx.RBAC.GetRole(id)
	if err != nil {
		return errors.New("角色不存在")
	}
	if role.Code == model.RoleCodeSuperAdmin && req.Code != "" && req.Code != model.RoleCodeSuperAdmin {
		return errors.New("不可修改超管编码")
	}
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Code != "" {
		role.Code = req.Code
	}
	role.Status = req.Status
	role.Remark = req.Remark
	return l.svcCtx.RBAC.UpdateRole(role)
}

func (l *RBACLogic) DeleteRole(id uint64) error {
	role, err := l.svcCtx.RBAC.GetRole(id)
	if err != nil {
		return errors.New("角色不存在")
	}
	if role.Code == model.RoleCodeSuperAdmin {
		return errors.New("不可删除超级管理员角色")
	}
	return l.svcCtx.RBAC.DeleteRole(id)
}

func (l *RBACLogic) GetRoleMenus(id uint64) ([]uint64, error) {
	return l.svcCtx.RBAC.ListRoleMenuIDs(id)
}

func (l *RBACLogic) AssignRoleMenus(id uint64, menuIDs []uint64) error {
	if _, err := l.svcCtx.RBAC.GetRole(id); err != nil {
		return errors.New("角色不存在")
	}
	return l.svcCtx.RBAC.ReplaceRoleMenus(id, menuIDs)
}

func (l *RBACLogic) ListUsers(page, pageSize int, mobile string) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return l.svcCtx.RBAC.ListUsers(page, pageSize, mobile)
}

func (l *RBACLogic) SetUserStatus(id uint64, status int) error {
	if status != 0 && status != 1 {
		return errors.New("status 无效")
	}
	return l.svcCtx.RBAC.UpdateUserStatus(id, status)
}

func (l *RBACLogic) GetUser(id uint64) (*model.User, error) {
	user, err := l.svcCtx.Repo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func (l *RBACLogic) UpdateUser(id uint64, req types.UserUpdateReq) error {
	if _, err := l.svcCtx.Repo.FindByID(id); err != nil {
		return errors.New("用户不存在")
	}
	if len(req.Mobile) != 11 {
		return errors.New("手机号无效")
	}
	if l.svcCtx.Repo.MobileTakenByOther(req.Mobile, id) {
		return errors.New("手机号已被占用")
	}
	if req.Gender < 0 || req.Gender > 2 {
		return errors.New("性别无效")
	}
	return l.svcCtx.Repo.UpdateProfile(id, req.Nickname, req.Avatar, req.Mobile, req.Gender)
}

func (l *RBACLogic) ResetUserPassword(id uint64, plain string) error {
	if plain == "" {
		return errors.New("密码不能为空")
	}
	if _, err := l.svcCtx.Repo.FindByID(id); err != nil {
		return errors.New("用户不存在")
	}
	return l.svcCtx.Repo.UpdatePassword(id, plain)
}

func (l *RBACLogic) ListAdmins(page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return l.svcCtx.RBAC.ListAdmins(page, pageSize)
}

func (l *RBACLogic) CreateAdmin(req types.AdminCreateReq) (*model.User, error) {
	if len(req.Mobile) != 11 || req.Password == "" {
		return nil, errors.New("手机号或密码无效")
	}
	user, err := l.svcCtx.Repo.CreateAdmin(req.Mobile, req.Password, req.Nickname)
	if err != nil {
		return nil, err
	}
	if len(req.RoleIDs) > 0 {
		if err := l.svcCtx.RBAC.ReplaceUserRoles(user.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (l *RBACLogic) AssignAdminRoles(userID uint64, roleIDs []uint64) error {
	if _, err := l.svcCtx.Repo.FindByID(userID); err != nil {
		return errors.New("用户不存在")
	}
	return l.svcCtx.RBAC.ReplaceUserRoles(userID, roleIDs)
}

func (l *RBACLogic) AdminRoleIDs(userID uint64) ([]uint64, error) {
	return l.svcCtx.RBAC.ListUserRoleIDs(userID)
}

func (l *RBACLogic) ResetAdminPassword(userID uint64, password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}
	return l.svcCtx.Repo.UpdatePassword(userID, password)
}

func (l *RBACLogic) ListConfigs() ([]model.SysConfig, error) {
	return l.svcCtx.RBAC.ListConfigs()
}

func (l *RBACLogic) SaveConfigs(items []types.ConfigItemReq) error {
	for _, it := range items {
		if it.ConfigKey == "" {
			continue
		}
		if err := l.svcCtx.RBAC.UpsertConfig(it.ConfigKey, it.ConfigValue, it.Remark); err != nil {
			return err
		}
	}
	return nil
}

// ResolveLoginRole 有平台角色则 JWT 使用 platform_admin
func ResolveLoginRole(rbac *svc.ServiceContext, user *model.User) string {
	role := user.Role
	if role == "" {
		role = jwt.RoleUser
	}
	if rbac.RBAC.HasPlatformRole(user.ID) {
		return jwt.RolePlatformAdmin
	}
	return role
}
