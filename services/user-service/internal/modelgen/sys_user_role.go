package modelgen

type SysUserRole struct {
	UserID uint64 `db:"user_id"`
	RoleID uint64 `db:"role_id"`
}

func (SysUserRole) TableName() string { return "sys_user_role" }
