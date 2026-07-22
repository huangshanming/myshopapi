package modelgen

type ShopUserRole struct {
	ShopID uint64 `db:"shop_id"`
	UserID uint64 `db:"user_id"`
	RoleID uint64 `db:"role_id"`
}

func (ShopUserRole) TableName() string { return "shop_user_roles" }
