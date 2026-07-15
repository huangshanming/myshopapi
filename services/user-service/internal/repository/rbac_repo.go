package repository

import (
	"mymall/services/user-service/internal/model"

	"gorm.io/gorm"
)

type RBACRepository struct {
	db *gorm.DB
}

func NewRBACRepository(db *gorm.DB) *RBACRepository {
	return &RBACRepository{db: db}
}

func (r *RBACRepository) DB() *gorm.DB { return r.db }

func (r *RBACRepository) HasPlatformRole(userID uint64) bool {
	var n int64
	_ = r.db.Model(&model.SysUserRole{}).Where("user_id = ?", userID).Count(&n).Error
	return n > 0
}

func (r *RBACRepository) IsSuperAdmin(userID uint64) bool {
	var n int64
	err := r.db.Table("sys_user_role ur").
		Joins("JOIN sys_role r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.code = ? AND r.status = 1", userID, model.RoleCodeSuperAdmin).
		Count(&n).Error
	return err == nil && n > 0
}

func (r *RBACRepository) ListRoleCodes(userID uint64) ([]string, error) {
	var codes []string
	err := r.db.Table("sys_user_role ur").
		Select("r.code").
		Joins("JOIN sys_role r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.status = 1", userID).
		Pluck("r.code", &codes).Error
	return codes, err
}

func (r *RBACRepository) ListUserPerms(userID uint64) ([]string, error) {
	if r.IsSuperAdmin(userID) {
		var all []string
		err := r.db.Model(&model.SysMenu{}).
			Where("status = 1 AND perms <> ''").
			Pluck("perms", &all).Error
		return all, err
	}
	var perms []string
	err := r.db.Table("sys_user_role ur").
		Select("DISTINCT m.perms").
		Joins("JOIN sys_role_menu rm ON rm.role_id = ur.role_id").
		Joins("JOIN sys_menu m ON m.id = rm.menu_id").
		Where("ur.user_id = ? AND m.status = 1 AND m.perms <> ''", userID).
		Pluck("m.perms", &perms).Error
	return perms, err
}

func (r *RBACRepository) ListUserMenus(userID uint64) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	q := r.db.Model(&model.SysMenu{}).Where("status = 1 AND type IN ?", []string{model.MenuTypeDir, model.MenuTypeMenu})
	if !r.IsSuperAdmin(userID) {
		q = r.db.Table("sys_menu m").
			Select("DISTINCT m.*").
			Joins("JOIN sys_role_menu rm ON rm.menu_id = m.id").
			Joins("JOIN sys_user_role ur ON ur.role_id = rm.role_id").
			Where("ur.user_id = ? AND m.status = 1 AND m.type IN ?", userID, []string{model.MenuTypeDir, model.MenuTypeMenu})
	}
	err := q.Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *RBACRepository) ListAllMenus() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.db.Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *RBACRepository) GetMenu(id uint64) (*model.SysMenu, error) {
	var m model.SysMenu
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *RBACRepository) CreateMenu(m *model.SysMenu) error {
	return r.db.Create(m).Error
}

func (r *RBACRepository) UpdateMenu(m *model.SysMenu) error {
	return r.db.Save(m).Error
}

func (r *RBACRepository) DeleteMenu(id uint64) error {
	var child int64
	_ = r.db.Model(&model.SysMenu{}).Where("parent_id = ?", id).Count(&child).Error
	if child > 0 {
		return gorm.ErrForeignKeyViolated
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id = ?", id).Delete(&model.SysRoleMenu{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.SysMenu{}, id).Error
	})
}

func (r *RBACRepository) ListRoles() ([]model.SysRole, error) {
	var roles []model.SysRole
	err := r.db.Order("id ASC").Find(&roles).Error
	return roles, err
}

func (r *RBACRepository) GetRole(id uint64) (*model.SysRole, error) {
	var role model.SysRole
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RBACRepository) CreateRole(role *model.SysRole) error {
	return r.db.Create(role).Error
}

func (r *RBACRepository) UpdateRole(role *model.SysRole) error {
	return r.db.Save(role).Error
}

func (r *RBACRepository) DeleteRole(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&model.SysRoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&model.SysUserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.SysRole{}, id).Error
	})
}

func (r *RBACRepository) ListRoleMenuIDs(roleID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.Model(&model.SysRoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &ids).Error
	return ids, err
}

func (r *RBACRepository) ReplaceRoleMenus(roleID uint64, menuIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.SysRoleMenu{}).Error; err != nil {
			return err
		}
		for _, mid := range menuIDs {
			if err := tx.Create(&model.SysRoleMenu{RoleID: roleID, MenuID: mid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RBACRepository) ListUsers(page, pageSize int, mobile string) ([]model.User, int64, error) {
	q := r.db.Model(&model.User{})
	if mobile != "" {
		q = q.Where("mobile LIKE ?", "%"+mobile+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.User
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *RBACRepository) UpdateUserStatus(id uint64, status int) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *RBACRepository) ListAdmins(page, pageSize int) ([]model.User, int64, error) {
	sub := r.db.Model(&model.SysUserRole{}).Select("DISTINCT user_id")
	q := r.db.Model(&model.User{}).Where("id IN (?) OR role = ?", sub, "platform_admin")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.User
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *RBACRepository) ListUserRoleIDs(userID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.Model(&model.SysUserRole{}).Where("user_id = ?", userID).Pluck("role_id", &ids).Error
	return ids, err
}

func (r *RBACRepository) ReplaceUserRoles(userID uint64, roleIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.SysUserRole{}).Error; err != nil {
			return err
		}
		for _, rid := range roleIDs {
			if err := tx.Create(&model.SysUserRole{UserID: userID, RoleID: rid}).Error; err != nil {
				return err
			}
		}
		role := "user"
		if len(roleIDs) > 0 {
			role = "platform_admin"
		}
		return tx.Model(&model.User{}).Where("id = ?", userID).Update("role", role).Error
	})
}

func (r *RBACRepository) ListConfigs() ([]model.SysConfig, error) {
	var list []model.SysConfig
	err := r.db.Order("id ASC").Find(&list).Error
	return list, err
}

func (r *RBACRepository) UpsertConfig(key, value, remark string) error {
	var c model.SysConfig
	err := r.db.Where("config_key = ?", key).First(&c).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(&model.SysConfig{ConfigKey: key, ConfigValue: value, Remark: remark}).Error
	}
	if err != nil {
		return err
	}
	c.ConfigValue = value
	if remark != "" {
		c.Remark = remark
	}
	return r.db.Save(&c).Error
}
