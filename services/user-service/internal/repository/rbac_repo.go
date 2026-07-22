package repository

import (
	"context"
	"errors"

	"mymall/services/user-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const menuColumns = "id, created_at, updated_at, parent_id, IFNULL(name,'') AS name, IFNULL(type,'') AS type, IFNULL(path,'') AS path, IFNULL(component,'') AS component, IFNULL(icon,'') AS icon, IFNULL(perms,'') AS perms, sort, visible, status"
const roleColumns = "id, created_at, updated_at, IFNULL(code,'') AS code, IFNULL(name,'') AS name, status, IFNULL(remark,'') AS remark"
const configColumns = "id, created_at, updated_at, IFNULL(config_key,'') AS config_key, IFNULL(config_value,'') AS config_value, IFNULL(remark,'') AS remark"

type RBACRepository struct {
	conn sqlx.SqlConn
}

func NewRBACRepository(conn sqlx.SqlConn) *RBACRepository {
	return &RBACRepository{conn: conn}
}

func (r *RBACRepository) HasPlatformRole(ctx context.Context, userID uint64) bool {
	n, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM sys_user_role WHERE user_id=?", userID,
	)
	return err == nil && n > 0
}

func (r *RBACRepository) IsSuperAdmin(ctx context.Context, userID uint64) bool {
	n, err := countQuery(ctx, r.conn,
		`SELECT COUNT(*) FROM sys_user_role ur
		 JOIN sys_role r ON r.id = ur.role_id
		 WHERE ur.user_id=? AND r.code=? AND r.status=1`,
		userID, model.RoleCodeSuperAdmin,
	)
	return err == nil && n > 0
}

func (r *RBACRepository) ListRoleCodes(ctx context.Context, userID uint64) ([]string, error) {
	type row struct {
		Code string `db:"code"`
	}
	var rows []row
	err := r.conn.QueryRowsPartialCtx(ctx, &rows,
		`SELECT r.code FROM sys_user_role ur
		 JOIN sys_role r ON r.id = ur.role_id
		 WHERE ur.user_id=? AND r.status=1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(rows))
	for _, item := range rows {
		codes = append(codes, item.Code)
	}
	return codes, nil
}

func (r *RBACRepository) ListUserPerms(ctx context.Context, userID uint64) ([]string, error) {
	if r.IsSuperAdmin(ctx, userID) {
		type row struct {
			Perms string `db:"perms"`
		}
		var rows []row
		err := r.conn.QueryRowsPartialCtx(ctx, &rows,
			"SELECT perms FROM sys_menu WHERE status=1 AND perms<>''",
		)
		if err != nil {
			return nil, err
		}
		all := make([]string, 0, len(rows))
		for _, item := range rows {
			all = append(all, item.Perms)
		}
		return all, nil
	}
	type row struct {
		Perms string `db:"perms"`
	}
	var rows []row
	err := r.conn.QueryRowsPartialCtx(ctx, &rows,
		`SELECT DISTINCT m.perms FROM sys_user_role ur
		 JOIN sys_role_menu rm ON rm.role_id = ur.role_id
		 JOIN sys_menu m ON m.id = rm.menu_id
		 WHERE ur.user_id=? AND m.status=1 AND m.perms<>''`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	perms := make([]string, 0, len(rows))
	for _, item := range rows {
		perms = append(perms, item.Perms)
	}
	return perms, nil
}

func (r *RBACRepository) ListUserMenus(ctx context.Context, userID uint64) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	if !r.IsSuperAdmin(ctx, userID) {
		err := r.conn.QueryRowsPartialCtx(ctx, &menus,
			`SELECT DISTINCT `+menuColumns+` FROM sys_menu m
			 JOIN sys_role_menu rm ON rm.menu_id = m.id
			 JOIN sys_user_role ur ON ur.role_id = rm.role_id
			 WHERE ur.user_id=? AND m.status=1 AND m.type IN (?, ?)
			 ORDER BY m.sort ASC, m.id ASC`,
			userID, model.MenuTypeDir, model.MenuTypeMenu,
		)
		return menus, err
	}
	err := r.conn.QueryRowsPartialCtx(ctx, &menus,
		"SELECT "+menuColumns+" FROM sys_menu WHERE status=1 AND type IN (?, ?) ORDER BY sort ASC, id ASC",
		model.MenuTypeDir, model.MenuTypeMenu,
	)
	return menus, err
}

func (r *RBACRepository) ListAllMenus(ctx context.Context) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.conn.QueryRowsPartialCtx(ctx, &menus,
		"SELECT "+menuColumns+" FROM sys_menu ORDER BY sort ASC, id ASC",
	)
	return menus, err
}

func (r *RBACRepository) GetMenu(ctx context.Context, id uint64) (*model.SysMenu, error) {
	var m model.SysMenu
	err := r.conn.QueryRowPartialCtx(ctx, &m,
		"SELECT "+menuColumns+" FROM sys_menu WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *RBACRepository) CreateMenu(ctx context.Context, m *model.SysMenu) error {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO sys_menu (parent_id, name, type, path, component, icon, perms, sort, visible, status)
		 VALUES (?,?,?,?,?,?,?,?,?,?)`,
		m.ParentID, m.Name, m.Type, m.Path, m.Component, m.Icon, m.Perms, m.Sort, m.Visible, m.Status,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

func (r *RBACRepository) UpdateMenu(ctx context.Context, m *model.SysMenu) error {
	_, err := r.conn.ExecCtx(ctx,
		`UPDATE sys_menu SET parent_id=?, name=?, type=?, path=?, component=?, icon=?, perms=?, sort=?, visible=?, status=?
		 WHERE id=?`,
		m.ParentID, m.Name, m.Type, m.Path, m.Component, m.Icon, m.Perms, m.Sort, m.Visible, m.Status, m.ID,
	)
	return err
}

func (r *RBACRepository) DeleteMenu(ctx context.Context, id uint64) error {
	child, _ := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM sys_menu WHERE parent_id=?", id,
	)
	if child > 0 {
		return ErrMenuHasChildren
	}
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, "DELETE FROM sys_role_menu WHERE menu_id=?", id); err != nil {
			return err
		}
		_, err := session.ExecCtx(ctx, "DELETE FROM sys_menu WHERE id=?", id)
		return err
	})
}

func (r *RBACRepository) ListRoles(ctx context.Context) ([]model.SysRole, error) {
	var roles []model.SysRole
	err := r.conn.QueryRowsPartialCtx(ctx, &roles,
		"SELECT "+roleColumns+" FROM sys_role ORDER BY id ASC",
	)
	return roles, err
}

func (r *RBACRepository) GetRole(ctx context.Context, id uint64) (*model.SysRole, error) {
	var role model.SysRole
	err := r.conn.QueryRowPartialCtx(ctx, &role,
		"SELECT "+roleColumns+" FROM sys_role WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RBACRepository) CreateRole(ctx context.Context, role *model.SysRole) error {
	res, err := r.conn.ExecCtx(ctx,
		"INSERT INTO sys_role (code, name, status, remark) VALUES (?,?,?,?)",
		role.Code, role.Name, role.Status, role.Remark,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	role.ID = id
	return nil
}

func (r *RBACRepository) UpdateRole(ctx context.Context, role *model.SysRole) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE sys_role SET code=?, name=?, status=?, remark=? WHERE id=?",
		role.Code, role.Name, role.Status, role.Remark, role.ID,
	)
	return err
}

func (r *RBACRepository) DeleteRole(ctx context.Context, id uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, "DELETE FROM sys_role_menu WHERE role_id=?", id); err != nil {
			return err
		}
		if _, err := session.ExecCtx(ctx, "DELETE FROM sys_user_role WHERE role_id=?", id); err != nil {
			return err
		}
		_, err := session.ExecCtx(ctx, "DELETE FROM sys_role WHERE id=?", id)
		return err
	})
}

func (r *RBACRepository) ListRoleMenuIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	type row struct {
		MenuID uint64 `db:"menu_id"`
	}
	var rows []row
	err := r.conn.QueryRowsPartialCtx(ctx, &rows,
		"SELECT menu_id FROM sys_role_menu WHERE role_id=?", roleID,
	)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(rows))
	for _, item := range rows {
		ids = append(ids, item.MenuID)
	}
	return ids, nil
}

func (r *RBACRepository) ReplaceRoleMenus(ctx context.Context, roleID uint64, menuIDs []uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, "DELETE FROM sys_role_menu WHERE role_id=?", roleID); err != nil {
			return err
		}
		for _, mid := range menuIDs {
			if _, err := session.ExecCtx(ctx,
				"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?,?)", roleID, mid,
			); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RBACRepository) ListUsers(ctx context.Context, page, pageSize int, mobile string) ([]model.User, int64, error) {
	where := "deleted_at IS NULL"
	args := make([]any, 0, 4)
	if mobile != "" {
		where += " AND mobile LIKE ?"
		args = append(args, "%"+mobile+"%")
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM users WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.User
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+userColumns+" FROM users WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *RBACRepository) UpdateUserStatus(ctx context.Context, id uint64, status int) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE users SET status=? WHERE id=? AND deleted_at IS NULL", status, id,
	)
	return err
}

func (r *RBACRepository) ListAdmins(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	where := "deleted_at IS NULL AND (id IN (SELECT DISTINCT user_id FROM sys_user_role) OR role=?)"
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM users WHERE "+where, "platform_admin",
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.User
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+userColumns+" FROM users WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		"platform_admin", pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *RBACRepository) ListUserRoleIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	type row struct {
		RoleID uint64 `db:"role_id"`
	}
	var rows []row
	err := r.conn.QueryRowsPartialCtx(ctx, &rows,
		"SELECT role_id FROM sys_user_role WHERE user_id=?", userID,
	)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(rows))
	for _, item := range rows {
		ids = append(ids, item.RoleID)
	}
	return ids, nil
}

func (r *RBACRepository) ReplaceUserRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, "DELETE FROM sys_user_role WHERE user_id=?", userID); err != nil {
			return err
		}
		for _, rid := range roleIDs {
			if _, err := session.ExecCtx(ctx,
				"INSERT INTO sys_user_role (user_id, role_id) VALUES (?,?)", userID, rid,
			); err != nil {
				return err
			}
		}
		role := "user"
		if len(roleIDs) > 0 {
			role = "platform_admin"
		}
		_, err := session.ExecCtx(ctx,
			"UPDATE users SET role=? WHERE id=? AND deleted_at IS NULL", role, userID,
		)
		return err
	})
}

func (r *RBACRepository) ListConfigs(ctx context.Context) ([]model.SysConfig, error) {
	var list []model.SysConfig
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+configColumns+" FROM sys_config ORDER BY id ASC",
	)
	return list, err
}

func (r *RBACRepository) UpsertConfig(ctx context.Context, key, value, remark string) error {
	var c model.SysConfig
	err := r.conn.QueryRowPartialCtx(ctx, &c,
		"SELECT "+configColumns+" FROM sys_config WHERE config_key=? LIMIT 1", key,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		_, err = r.conn.ExecCtx(ctx,
			"INSERT INTO sys_config (config_key, config_value, remark) VALUES (?,?,?)",
			key, value, remark,
		)
		return err
	}
	if err != nil {
		return err
	}
	if remark != "" {
		c.Remark = remark
	}
	_, err = r.conn.ExecCtx(ctx,
		"UPDATE sys_config SET config_value=?, remark=? WHERE id=?",
		value, c.Remark, c.ID,
	)
	return err
}
