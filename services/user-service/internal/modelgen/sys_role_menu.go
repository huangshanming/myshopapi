package modelgen

type SysRoleMenu struct {
	RoleID uint64 `db:"role_id"`
	MenuID uint64 `db:"menu_id"`
}

func (SysRoleMenu) TableName() string { return "sys_role_menu" }
