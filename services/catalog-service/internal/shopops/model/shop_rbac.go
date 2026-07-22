package model

import "mymall/common"

type ShopRole struct {
	ID        uint64           `db:"id" json:"id"`
	ShopID    uint64           `db:"shop_id" json:"shop_id"`
	Code      string           `db:"code" json:"code"`
	Name      string           `db:"name" json:"name"`
	Status    int              `db:"status" json:"status"`
	Remark    string           `db:"remark" json:"remark"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (ShopRole) TableName() string { return "shop_roles" }

type ShopMenu struct {
	ID        uint64           `db:"id" json:"id"`
	ParentID  uint64           `db:"parent_id" json:"parent_id"`
	Name      string           `db:"name" json:"name"`
	Type      string           `db:"type" json:"type"`
	Path      string           `db:"path" json:"path"`
	Component string           `db:"component" json:"component"`
	Icon      string           `db:"icon" json:"icon"`
	Perms     string           `db:"perms" json:"perms"`
	Sort      int              `db:"sort" json:"sort"`
	Visible   int              `db:"visible" json:"visible"`
	Status    int              `db:"status" json:"status"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (ShopMenu) TableName() string { return "shop_menus" }

type ShopRoleMenu struct {
	RoleID uint64 `db:"role_id" json:"role_id"`
	MenuID uint64 `db:"menu_id" json:"menu_id"`
}

func (ShopRoleMenu) TableName() string { return "shop_role_menus" }

type ShopUserRole struct {
	ShopID uint64 `db:"shop_id" json:"shop_id"`
	UserID uint64 `db:"user_id" json:"user_id"`
	RoleID uint64 `db:"role_id" json:"role_id"`
}

func (ShopUserRole) TableName() string { return "shop_user_roles" }
