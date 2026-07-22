package modelgen

type ShopRoleMenu struct {
	RoleID uint64 `db:"role_id"`
	MenuID uint64 `db:"menu_id"`
}

func (ShopRoleMenu) TableName() string { return "shop_role_menus" }
