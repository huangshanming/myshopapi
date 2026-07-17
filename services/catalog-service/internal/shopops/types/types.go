package types

type ShopRoleReq struct {
	Code    string   `json:"code"`
	Name    string   `json:"name"`
	Remark  string   `json:"remark"`
	MenuIDs []uint64 `json:"menu_ids"`
}

type ShopStaffReq struct {
	Mobile   string `json:"mobile"`
	RoleID   uint64 `json:"role_id"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Mode     string `json:"mode"` // bind=绑定已有账号（默认） create=新建账号并绑定
}
