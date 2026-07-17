package model

import "mymall/common"

type ShopRole struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID    uint64           `gorm:"column:shop_id;not null" json:"shop_id"`
	Code      string           `gorm:"column:code;type:varchar(64)" json:"code"`
	Name      string           `gorm:"column:name;type:varchar(100)" json:"name"`
	Status    int              `gorm:"column:status" json:"status"`
	Remark    string           `gorm:"column:remark;type:varchar(255)" json:"remark"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ShopRole) TableName() string { return "shop_roles" }

type ShopMenu struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentID  uint64           `gorm:"column:parent_id" json:"parent_id"`
	Name      string           `gorm:"column:name" json:"name"`
	Type      string           `gorm:"column:type" json:"type"`
	Path      string           `gorm:"column:path" json:"path"`
	Component string           `gorm:"column:component" json:"component"`
	Icon      string           `gorm:"column:icon" json:"icon"`
	Perms     string           `gorm:"column:perms" json:"perms"`
	Sort      int              `gorm:"column:sort" json:"sort"`
	Visible   int              `gorm:"column:visible" json:"visible"`
	Status    int              `gorm:"column:status" json:"status"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ShopMenu) TableName() string { return "shop_menus" }

type ShopRoleMenu struct {
	RoleID uint64 `gorm:"column:role_id;primaryKey" json:"role_id"`
	MenuID uint64 `gorm:"column:menu_id;primaryKey" json:"menu_id"`
}

func (ShopRoleMenu) TableName() string { return "shop_role_menus" }

type ShopUserRole struct {
	ShopID uint64 `gorm:"column:shop_id;primaryKey" json:"shop_id"`
	UserID uint64 `gorm:"column:user_id;primaryKey" json:"user_id"`
	RoleID uint64 `gorm:"column:role_id;primaryKey" json:"role_id"`
}

func (ShopUserRole) TableName() string { return "shop_user_roles" }
