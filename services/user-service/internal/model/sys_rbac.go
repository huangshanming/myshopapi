package model

const (
	MenuTypeDir    = "dir"
	MenuTypeMenu   = "menu"
	MenuTypeButton = "button"

	RoleCodeSuperAdmin = "super_admin"
)

type SysMenu struct {
	BaseModel
	ParentID  uint64 `gorm:"column:parent_id;not null;default:0;index" json:"parent_id"`
	Name      string `gorm:"column:name;type:varchar(64);not null" json:"name"`
	Type      string `gorm:"column:type;type:varchar(16);not null;default:menu" json:"type"`
	Path      string `gorm:"column:path;type:varchar(128);not null;default:''" json:"path"`
	Component string `gorm:"column:component;type:varchar(128);not null;default:''" json:"component"`
	Icon      string `gorm:"column:icon;type:varchar(64);not null;default:''" json:"icon"`
	Perms     string `gorm:"column:perms;type:varchar(128);not null;default:''" json:"perms"`
	Sort      int    `gorm:"column:sort;not null;default:0" json:"sort"`
	Visible   int    `gorm:"column:visible;type:tinyint;not null;default:1" json:"visible"`
	Status    int    `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
}

func (SysMenu) TableName() string { return "sys_menu" }

type SysRole struct {
	BaseModel
	Code   string `gorm:"column:code;type:varchar(64);not null;uniqueIndex" json:"code"`
	Name   string `gorm:"column:name;type:varchar(64);not null" json:"name"`
	Status int    `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	Remark string `gorm:"column:remark;type:varchar(255);not null;default:''" json:"remark"`
}

func (SysRole) TableName() string { return "sys_role" }

type SysRoleMenu struct {
	RoleID uint64 `gorm:"column:role_id;primaryKey" json:"role_id"`
	MenuID uint64 `gorm:"column:menu_id;primaryKey" json:"menu_id"`
}

func (SysRoleMenu) TableName() string { return "sys_role_menu" }

type SysUserRole struct {
	UserID uint64 `gorm:"column:user_id;primaryKey" json:"user_id"`
	RoleID uint64 `gorm:"column:role_id;primaryKey" json:"role_id"`
}

func (SysUserRole) TableName() string { return "sys_user_role" }

type SysConfig struct {
	BaseModel
	ConfigKey   string `gorm:"column:config_key;type:varchar(64);not null;uniqueIndex" json:"config_key"`
	ConfigValue string `gorm:"column:config_value;type:varchar(512);not null;default:''" json:"config_value"`
	Remark      string `gorm:"column:remark;type:varchar(255);not null;default:''" json:"remark"`
}

func (SysConfig) TableName() string { return "sys_config" }
