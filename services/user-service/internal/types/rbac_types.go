package types

type MenuReq struct {
	ParentID  uint64 `json:"parent_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Icon      string `json:"icon"`
	Perms     string `json:"perms"`
	Sort      int    `json:"sort"`
	Visible   int    `json:"visible"`
	Status    int    `json:"status"`
}

type RoleReq struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

type RoleMenusReq struct {
	MenuIDs []uint64 `json:"menu_ids"`
}

type UserStatusReq struct {
	Status int `json:"status"`
}

type UserUpdateReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Mobile   string `json:"mobile"`
}

type UserResetPwdReq struct {
	Password string `json:"password"`
}

type AdminCreateReq struct {
	Mobile   string   `json:"mobile"`
	Password string   `json:"password"`
	Nickname string   `json:"nickname"`
	RoleIDs  []uint64 `json:"role_ids"`
}

type AdminRolesReq struct {
	RoleIDs []uint64 `json:"role_ids"`
}

type AdminResetPwdReq struct {
	Password string `json:"password"`
}

type ConfigItemReq struct {
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Remark      string `json:"remark"`
}

type ConfigBatchReq struct {
	Items []ConfigItemReq `json:"items"`
}

type PageListResp struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

type MenuTreeNode struct {
	ID        uint64          `json:"id"`
	ParentID  uint64          `json:"parent_id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Path      string          `json:"path"`
	Component string          `json:"component"`
	Icon      string          `json:"icon"`
	Perms     string          `json:"perms"`
	Sort      int             `json:"sort"`
	Visible   int             `json:"visible"`
	Status    int             `json:"status"`
	Children  []*MenuTreeNode `json:"children,omitempty"`
}

type AuthMeResp struct {
	Roles    []string        `json:"roles"`
	Perms    []string        `json:"perms"`
	MenuTree []*MenuTreeNode `json:"menu_tree"`
}
