package types

// ApplyReq 商家入驻申请
type ApplyReq struct {
	ShopName          string `json:"shop_name"`
	ContactName       string `json:"contact_name"`
	ContactPhone      string `json:"contact_phone"`
	Description       string `json:"description"`
	Category          string `json:"category"`
	Province          string `json:"province"`
	City              string `json:"city"`
	District          string `json:"district"`
	Address           string `json:"address"`
	BusinessLicenseNo string `json:"business_license_no"`
	LegalPerson       string `json:"legal_person"`
	LicenseImage      string `json:"license_image"`
	StorefrontImage   string `json:"storefront_image"`
}

// RejectReq 拒绝申请 / 禁用店铺
type RejectReq struct {
	Reason string `json:"reason"`
}

// UpdateShopReq 商家侧更新店铺资料（展示类）
type UpdateShopReq struct {
	Name            string `json:"name"`
	Logo            string `json:"logo"`
	ContactName     string `json:"contact_name"`
	ContactPhone    string `json:"contact_phone"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	Province        string `json:"province"`
	City            string `json:"city"`
	District        string `json:"district"`
	Address         string `json:"address"`
	StorefrontImage string `json:"storefront_image"`
}

// AdminCreateShopReq 平台后台直接开店
type AdminCreateShopReq struct {
	Name              string `json:"name"`
	Logo              string `json:"logo"`
	ContactName       string `json:"contact_name"`
	ContactPhone      string `json:"contact_phone"`
	Description       string `json:"description"`
	Category          string `json:"category"`
	Province          string `json:"province"`
	City              string `json:"city"`
	District          string `json:"district"`
	Address           string `json:"address"`
	BusinessLicenseNo string `json:"business_license_no"`
	LegalPerson       string `json:"legal_person"`
	LicenseImage      string `json:"license_image"`
	StorefrontImage   string `json:"storefront_image"`
	OwnerMobile       string `json:"owner_mobile"`
	OwnerPassword     string `json:"owner_password"`
	OwnerNickname     string `json:"owner_nickname"`
}

// AdminUpdateShopReq 平台更新店铺资料（含入驻档）
type AdminUpdateShopReq struct {
	Name              string `json:"name"`
	Logo              string `json:"logo"`
	ContactName       string `json:"contact_name"`
	ContactPhone      string `json:"contact_phone"`
	Description       string `json:"description"`
	Category          string `json:"category"`
	Province          string `json:"province"`
	City              string `json:"city"`
	District          string `json:"district"`
	Address           string `json:"address"`
	BusinessLicenseNo string `json:"business_license_no"`
	LegalPerson       string `json:"legal_person"`
	LicenseImage      string `json:"license_image"`
	StorefrontImage   string `json:"storefront_image"`
}

// OwnerPasswordReq 重置店主密码
type OwnerPasswordReq struct {
	Password string `json:"password"`
}

// PageListResp 分页列表
type PageListResp struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
