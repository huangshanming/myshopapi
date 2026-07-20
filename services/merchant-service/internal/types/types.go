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

// WalletAdjustReq 平台调账
// Field: balance（可用余额，默认）/ deposit（保证金）/ frozen_balance（冻结余额）
type WalletAdjustReq struct {
	Field  string  `json:"field"`
	Amount float64 `json:"amount"`
	Remark string  `json:"remark"`
}

// SeckillRuleReq 更新秒杀规则
type SeckillRuleReq struct {
	DurationHours     int     `json:"duration_hours"`
	ApplyFee          float64 `json:"apply_fee"`
	MaxEntriesPerShop int     `json:"max_entries_per_shop"`
	Status            int8    `json:"status"`
}

// SeckillApplyReq 商家报名秒杀
type SeckillApplyReq struct {
	SessionID    uint64  `json:"session_id"`
	ProductID    uint64  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	OriginPrice  float64 `json:"origin_price"`
	SeckillPrice float64 `json:"seckill_price"`
	SeckillStock int     `json:"seckill_stock"`
	AutoRenew    *int8   `json:"auto_renew"` // 1=到期自动续费
}

// SeckillAutoRenewReq 开关自动续费
type SeckillAutoRenewReq struct {
	AutoRenew int8 `json:"auto_renew"`
}

// SeckillConsumeReq order-service 扣秒杀库存
type SeckillConsumeReq struct {
	EntryID   uint64 `json:"entry_id"`
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// SeckillRestoreReq 回补秒杀库存
type SeckillRestoreReq struct {
	EntryID  uint64 `json:"entry_id"`
	Quantity int    `json:"quantity"`
}
