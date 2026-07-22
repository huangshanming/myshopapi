package model

const (
	MenuTypeDir    = "dir"
	MenuTypeMenu   = "menu"
	MenuTypeButton = "button"

	RoleCodeSuperAdmin = "super_admin"
)

type SysMenu struct {
	BaseModel
	ParentID  uint64 `gorm:"column:parent_id;not null;default:0;index" db:"parent_id" json:"parent_id"`
	Name      string `gorm:"column:name;type:varchar(64);not null" db:"name" json:"name"`
	Type      string `gorm:"column:type;type:varchar(16);not null;default:menu" db:"type" json:"type"`
	Path      string `gorm:"column:path;type:varchar(128);not null;default:''" db:"path" json:"path"`
	Component string `gorm:"column:component;type:varchar(128);not null;default:''" db:"component" json:"component"`
	Icon      string `gorm:"column:icon;type:varchar(64);not null;default:''" db:"icon" json:"icon"`
	Perms     string `gorm:"column:perms;type:varchar(128);not null;default:''" db:"perms" json:"perms"`
	Sort      int    `gorm:"column:sort;not null;default:0" db:"sort" json:"sort"`
	Visible   int    `gorm:"column:visible;type:tinyint;not null;default:1" db:"visible" json:"visible"`
	Status    int    `gorm:"column:status;type:tinyint;not null;default:1" db:"status" json:"status"`
}

func (SysMenu) TableName() string { return "sys_menu" }

type SysRole struct {
	BaseModel
	Code   string `gorm:"column:code;type:varchar(64);not null;uniqueIndex" db:"code" json:"code"`
	Name   string `gorm:"column:name;type:varchar(64);not null" db:"name" json:"name"`
	Status int    `gorm:"column:status;type:tinyint;not null;default:1" db:"status" json:"status"`
	Remark string `gorm:"column:remark;type:varchar(255);not null;default:''" db:"remark" json:"remark"`
}

func (SysRole) TableName() string { return "sys_role" }

type SysRoleMenu struct {
	RoleID uint64 `gorm:"column:role_id;primaryKey" db:"role_id" json:"role_id"`
	MenuID uint64 `gorm:"column:menu_id;primaryKey" db:"menu_id" json:"menu_id"`
}

func (SysRoleMenu) TableName() string { return "sys_role_menu" }

type SysUserRole struct {
	UserID uint64 `gorm:"column:user_id;primaryKey" db:"user_id" json:"user_id"`
	RoleID uint64 `gorm:"column:role_id;primaryKey" db:"role_id" json:"role_id"`
}

func (SysUserRole) TableName() string { return "sys_user_role" }

type SysConfig struct {
	BaseModel
	ConfigKey   string `gorm:"column:config_key;type:varchar(64);not null;uniqueIndex" db:"config_key" json:"config_key"`
	ConfigValue string `gorm:"column:config_value;type:varchar(512);not null;default:''" db:"config_value" json:"config_value"`
	Remark      string `gorm:"column:remark;type:varchar(255);not null;default:''" db:"remark" json:"remark"`
}

func (SysConfig) TableName() string { return "sys_config" }
