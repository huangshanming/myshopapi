package types

// ApplyReq 商家入驻申请
type ApplyReq struct {
	ShopName     string `json:"shop_name"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	Description  string `json:"description"`
}

// RejectReq 拒绝申请 / 禁用店铺
type RejectReq struct {
	Reason string `json:"reason"`
}

// UpdateShopReq 更新店铺资料
type UpdateShopReq struct {
	Name         string `json:"name"`
	Logo         string `json:"logo"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	Description  string `json:"description"`
}

// PageListResp 分页列表
type PageListResp struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
