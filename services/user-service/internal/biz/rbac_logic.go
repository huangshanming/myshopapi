package biz

import (
	"context"
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

func (l *RBACLogic) IsSuperAdmin(ctx context.Context, userID uint64) bool {
	return l.svcCtx.RBAC.IsSuperAdmin(ctx, userID)
}

func (l *RBACLogic) HasPerm(ctx context.Context, userID uint64, code string) bool {
	perms, err := l.svcCtx.RBAC.ListUserPerms(ctx, userID)
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

func (l *RBACLogic) AuthMe(ctx context.Context, userID uint64) (*types.AuthMeResp, error) {
	roles, err := l.svcCtx.RBAC.ListRoleCodes(ctx, userID)
	if err != nil {
		return nil, err
	}
	perms, err := l.svcCtx.RBAC.ListUserPerms(ctx, userID)
	if err != nil {
		return nil, err
	}
	menus, err := l.svcCtx.RBAC.ListUserMenus(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &types.AuthMeResp{
		Roles:    roles,
		Perms:    perms,
		MenuTree: BuildMenuTree(menus, true),
	}, nil
}

func (l *RBACLogic) MenuTreeAll(ctx context.Context) ([]*types.MenuTreeNode, error) {
	menus, err := l.svcCtx.RBAC.ListAllMenus(ctx)
	if err != nil {
		return nil, err
	}
	return BuildMenuTree(menus, false), nil
}

func (l *RBACLogic) CreateMenu(ctx context.Context, req types.MenuReq) (*model.SysMenu, error) {
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
	if err := l.svcCtx.RBAC.CreateMenu(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (l *RBACLogic) UpdateMenu(ctx context.Context, id uint64, req types.MenuReq) error {
	m, err := l.svcCtx.RBAC.GetMenu(ctx, id)
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
	return l.svcCtx.RBAC.UpdateMenu(ctx, m)
}

func (l *RBACLogic) DeleteMenu(ctx context.Context, id uint64) error {
	if err := l.svcCtx.RBAC.DeleteMenu(ctx, id); err != nil {
		return errors.New("请先删除子菜单")
	}
	return nil
}

func (l *RBACLogic) ListRoles(ctx context.Context) ([]model.SysRole, error) {
	return l.svcCtx.RBAC.ListRoles(ctx)
}

func (l *RBACLogic) CreateRole(ctx context.Context, req types.RoleReq) (*model.SysRole, error) {
	if req.Code == "" || req.Name == "" {
		return nil, errors.New("编码和名称不能为空")
	}
	role := &model.SysRole{Code: req.Code, Name: req.Name, Status: req.Status, Remark: req.Remark}
	if role.Status == 0 {
		role.Status = 1
	}
	if err := l.svcCtx.RBAC.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (l *RBACLogic) UpdateRole(ctx context.Context, id uint64, req types.RoleReq) error {
	role, err := l.svcCtx.RBAC.GetRole(ctx, id)
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
	return l.svcCtx.RBAC.UpdateRole(ctx, role)
}

func (l *RBACLogic) DeleteRole(ctx context.Context, id uint64) error {
	role, err := l.svcCtx.RBAC.GetRole(ctx, id)
	if err != nil {
		return errors.New("角色不存在")
	}
	if role.Code == model.RoleCodeSuperAdmin {
		return errors.New("不可删除超级管理员角色")
	}
	return l.svcCtx.RBAC.DeleteRole(ctx, id)
}

func (l *RBACLogic) GetRoleMenus(ctx context.Context, id uint64) ([]uint64, error) {
	return l.svcCtx.RBAC.ListRoleMenuIDs(ctx, id)
}

func (l *RBACLogic) AssignRoleMenus(ctx context.Context, id uint64, menuIDs []uint64) error {
	if _, err := l.svcCtx.RBAC.GetRole(ctx, id); err != nil {
		return errors.New("角色不存在")
	}
	return l.svcCtx.RBAC.ReplaceRoleMenus(ctx, id, menuIDs)
}

func (l *RBACLogic) ListUsers(ctx context.Context, page, pageSize int, mobile string) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return l.svcCtx.RBAC.ListUsers(ctx, page, pageSize, mobile)
}

func (l *RBACLogic) SetUserStatus(ctx context.Context, id uint64, status int) error {
	if status != 0 && status != 1 {
		return errors.New("status 无效")
	}
	return l.svcCtx.RBAC.UpdateUserStatus(ctx, id, status)
}

func (l *RBACLogic) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	user, err := l.svcCtx.Repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func (l *RBACLogic) UpdateUser(ctx context.Context, id uint64, req types.UserUpdateReq) error {
	if _, err := l.svcCtx.Repo.FindByID(ctx, id); err != nil {
		return errors.New("用户不存在")
	}
	if len(req.Mobile) != 11 {
		return errors.New("手机号无效")
	}
	if l.svcCtx.Repo.MobileTakenByOther(ctx, req.Mobile, id) {
		return errors.New("手机号已被占用")
	}
	if req.Gender < 0 || req.Gender > 2 {
		return errors.New("性别无效")
	}
	return l.svcCtx.Repo.UpdateProfile(ctx, id, req.Nickname, req.Avatar, req.Mobile, req.Gender)
}

func (l *RBACLogic) ResetUserPassword(ctx context.Context, id uint64, plain string) error {
	if plain == "" {
		return errors.New("密码不能为空")
	}
	if _, err := l.svcCtx.Repo.FindByID(ctx, id); err != nil {
		return errors.New("用户不存在")
	}
	return l.svcCtx.Repo.UpdatePassword(ctx, id, plain)
}

// GenerateUserToken 为指定用户签发与登录一致的 JWT（并发压测 / 调试用）。
func (l *RBACLogic) GenerateUserToken(ctx context.Context, id uint64) (map[string]interface{}, error) {
	user, err := l.svcCtx.Repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.Status != 1 {
		return nil, errors.New("用户已禁用")
	}
	role := ResolveLoginRole(l.svcCtx, user)
	var shopID uint64
	if jwt.IsMerchant(role) {
		shopID = l.svcCtx.Repo.FirstShopID(ctx, user.ID)
	}
	token, err := jwt.GenerateTokenWithShop(user.ID, role, shopID, l.svcCtx.JWT)
	if err != nil {
		return nil, err
	}
	expireHours := l.svcCtx.JWT.ExpireHours
	if expireHours <= 0 {
		expireHours = 24
	}
	return map[string]interface{}{
		"token":        token,
		"user_id":      user.ID,
		"mobile":       user.Mobile,
		"nickname":     user.Nickname,
		"role":         role,
		"shop_id":      shopID,
		"expire_hours": expireHours,
	}, nil
}

func (l *RBACLogic) ListAdmins(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return l.svcCtx.RBAC.ListAdmins(ctx, page, pageSize)
}

func (l *RBACLogic) CreateAdmin(ctx context.Context, req types.AdminCreateReq) (*model.User, error) {
	if len(req.Mobile) != 11 || req.Password == "" {
		return nil, errors.New("手机号或密码无效")
	}
	user, err := l.svcCtx.Repo.CreateAdmin(ctx, req.Mobile, req.Password, req.Nickname)
	if err != nil {
		return nil, err
	}
	if len(req.RoleIDs) > 0 {
		if err := l.svcCtx.RBAC.ReplaceUserRoles(ctx, user.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (l *RBACLogic) AssignAdminRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	if _, err := l.svcCtx.Repo.FindByID(ctx, userID); err != nil {
		return errors.New("用户不存在")
	}
	return l.svcCtx.RBAC.ReplaceUserRoles(ctx, userID, roleIDs)
}

func (l *RBACLogic) AdminRoleIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	return l.svcCtx.RBAC.ListUserRoleIDs(ctx, userID)
}

func (l *RBACLogic) ResetAdminPassword(ctx context.Context, userID uint64, password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}
	return l.svcCtx.Repo.UpdatePassword(ctx, userID, password)
}

func (l *RBACLogic) ListConfigs(ctx context.Context) ([]model.SysConfig, error) {
	return l.svcCtx.RBAC.ListConfigs(ctx)
}

func (l *RBACLogic) GetConfigValue(ctx context.Context, key string) (string, error) {
	return l.svcCtx.RBAC.GetConfigValue(ctx, key)
}

func (l *RBACLogic) SaveConfigs(ctx context.Context, items []types.ConfigItemReq) error {
	for _, it := range items {
		if it.ConfigKey == "" {
			continue
		}
		if err := l.svcCtx.RBAC.UpsertConfig(ctx, it.ConfigKey, it.ConfigValue, it.Remark); err != nil {
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
	if rbac.RBAC.HasPlatformRole(context.Background(), user.ID) {
		return jwt.RolePlatformAdmin
	}
	return role
}
