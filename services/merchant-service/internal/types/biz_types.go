package types

import "mymall/common"

// Hand-maintained API/biz DTOs (kept in sync with api/merchant.api).

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
	Name            string   `json:"name"`
	Logo            string   `json:"logo"`
	ContactName     string   `json:"contact_name"`
	ContactPhone    string   `json:"contact_phone"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	Province        string   `json:"province"`
	City            string   `json:"city"`
	District        string   `json:"district"`
	Address         string   `json:"address"`
	StorefrontImage string   `json:"storefront_image"`
	Latitude        float64  `json:"latitude"`
	Longitude       float64  `json:"longitude"`
	LocalEnabled    int      `json:"local_enabled"`
	Images          []string `json:"images,optional"`
}

// UpdateShopBodyReq path+body
type UpdateShopBodyReq struct {
	Id              uint64   `path:"id"`
	Name            string   `json:"name"`
	Logo            string   `json:"logo"`
	ContactName     string   `json:"contact_name"`
	ContactPhone    string   `json:"contact_phone"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	Province        string   `json:"province"`
	City            string   `json:"city"`
	District        string   `json:"district"`
	Address         string   `json:"address"`
	StorefrontImage string   `json:"storefront_image"`
	Latitude        float64  `json:"latitude,optional"`
	Longitude       float64  `json:"longitude,optional"`
	LocalEnabled    int      `json:"local_enabled,optional"`
	Images          []string `json:"images,optional"`
}

func (r *UpdateShopBodyReq) ToUpdateShopReq() UpdateShopReq {
	return UpdateShopReq{
		Name: r.Name, Logo: r.Logo, ContactName: r.ContactName, ContactPhone: r.ContactPhone,
		Description: r.Description, Category: r.Category, Province: r.Province, City: r.City,
		District: r.District, Address: r.Address, StorefrontImage: r.StorefrontImage,
		Latitude: r.Latitude, Longitude: r.Longitude, LocalEnabled: r.LocalEnabled, Images: r.Images,
	}
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

// AdminUpdateShopBodyReq path+body
type AdminUpdateShopBodyReq struct {
	Id                uint64 `path:"id"`
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

func (r *AdminUpdateShopBodyReq) ToAdminUpdateShopReq() AdminUpdateShopReq {
	return AdminUpdateShopReq{
		Name: r.Name, Logo: r.Logo, ContactName: r.ContactName, ContactPhone: r.ContactPhone,
		Description: r.Description, Category: r.Category, Province: r.Province, City: r.City,
		District: r.District, Address: r.Address, BusinessLicenseNo: r.BusinessLicenseNo,
		LegalPerson: r.LegalPerson, LicenseImage: r.LicenseImage, StorefrontImage: r.StorefrontImage,
	}
}

// OwnerPasswordReq 重置店主密码
type OwnerPasswordReq struct {
	Password string `json:"password"`
}

// OwnerPasswordBodyReq path+body
type OwnerPasswordBodyReq struct {
	Id       uint64 `path:"id"`
	Password string `json:"password"`
}

// RejectBodyReq path+body（拒绝申请 / 禁用店铺）
type RejectBodyReq struct {
	Id     uint64 `path:"id"`
	Reason string `json:"reason"`
}

// WalletAdjustReq 平台调账
type WalletAdjustReq struct {
	Field  string  `json:"field"`
	Amount float64 `json:"amount"`
	Remark string  `json:"remark"`
}

// WalletAdjustBodyReq path+body
type WalletAdjustBodyReq struct {
	Id     uint64  `path:"id"`
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
	AutoRenew    *int8   `json:"auto_renew,optional"`
}

// SeckillAutoRenewReq 开关自动续费
type SeckillAutoRenewReq struct {
	AutoRenew int8 `json:"auto_renew"`
}

// SeckillAutoRenewBodyReq path+body
type SeckillAutoRenewBodyReq struct {
	Id        uint64 `path:"id"`
	AutoRenew int8   `json:"auto_renew"`
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

// BuySlotReq 商家购买首页坑位
type BuySlotReq struct {
	PackageID uint64 `json:"package_id"`
	TargetID  uint64 `json:"target_id,optional"`
}

// GrantSlotReq 平台赠送坑位
type GrantSlotReq struct {
	ShopID    uint64 `json:"shop_id"`
	PackageID uint64 `json:"package_id"`
	TargetID  uint64 `json:"target_id,optional"`
}

// ThemeBuyReq 商家购买主题坑
type ThemeBuyReq struct {
	ThemeSlotID uint64 `json:"theme_slot_id"`
	PackageID   uint64 `json:"package_id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle,optional"`
	CoverURL    string `json:"cover_url"`
	LinkType    string `json:"link_type"`
	LinkID      uint64 `json:"link_id,optional"`
}

// ThemeGrantReq 平台赠送主题坑
type ThemeGrantReq struct {
	ShopID      uint64 `json:"shop_id"`
	ThemeSlotID uint64 `json:"theme_slot_id"`
	PackageID   uint64 `json:"package_id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle,optional"`
	CoverURL    string `json:"cover_url"`
	LinkType    string `json:"link_type"`
	LinkID      uint64 `json:"link_id,optional"`
}

// SlotPackageSaveReq 坑位套餐
type SlotPackageSaveReq struct {
	SlotType     string  `json:"slot_type"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	DurationDays int     `json:"duration_days"`
	Status       string  `json:"status,optional"`
	Sort         int     `json:"sort,optional"`
	Remark       string  `json:"remark,optional"`
}

// SlotPackageUpdateBodyReq path+body
type SlotPackageUpdateBodyReq struct {
	Id           uint64  `path:"id"`
	SlotType     string  `json:"slot_type"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	DurationDays int     `json:"duration_days"`
	Status       string  `json:"status,optional"`
	Sort         int     `json:"sort,optional"`
	Remark       string  `json:"remark,optional"`
}

// SlotSettingItem 首页坑位展示上限
type SlotSettingItem struct {
	SlotType  string `json:"slot_type"`
	HomeLimit int    `json:"home_limit"`
}

// UpdateSlotSettingsReq 批量更新坑位设置
type UpdateSlotSettingsReq struct {
	Items []SlotSettingItem `json:"items"`
}

// ThemePackageSaveReq 主题套餐
type ThemePackageSaveReq struct {
	ThemeSlotID  uint64  `json:"theme_slot_id,optional"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	DurationDays int     `json:"duration_days"`
	Status       string  `json:"status,optional"`
	Sort         int     `json:"sort,optional"`
	Remark       string  `json:"remark,optional"`
}

// ThemePackageUpdateBodyReq path+body
type ThemePackageUpdateBodyReq struct {
	Id           uint64  `path:"id"`
	ThemeSlotID  uint64  `json:"theme_slot_id,optional"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	DurationDays int     `json:"duration_days"`
	Status       string  `json:"status,optional"`
	Sort         int     `json:"sort,optional"`
	Remark       string  `json:"remark,optional"`
}

// ThemeSlotUpdateBodyReq 更新主题坑位配置
type ThemeSlotUpdateBodyReq struct {
	Id              uint64 `path:"id"`
	Name            string `json:"name,optional"`
	Desc            string `json:"desc,optional"`
	CoverURL        string `json:"cover_url,optional"`
	DefaultLinkType string `json:"default_link_type,optional"`
	DefaultLinkID   uint64 `json:"default_link_id,optional"`
	Status          string `json:"status,optional"`
	Sort            int    `json:"sort,optional"`
	Position        string `json:"position,optional"`
}

// CouponScopeItem 优惠券适用范围
type CouponScopeItem struct {
	RefType string `json:"ref_type"`
	RefID   uint64 `json:"ref_id"`
}

// CouponSaveReq 创建/更新优惠券
type CouponSaveReq struct {
	Name              string             `json:"name"`
	CouponType        string             `json:"coupon_type"`
	ThresholdAmount   float64            `json:"threshold_amount,optional"`
	DiscountAmount    float64            `json:"discount_amount,optional"`
	DiscountRate      float64            `json:"discount_rate,optional"`
	MaxDiscountAmount float64            `json:"max_discount_amount,optional"`
	ScopeType         string             `json:"scope_type,optional"`
	TotalCount        int                `json:"total_count,optional"`
	PerUserLimit      int                `json:"per_user_limit,optional"`
	ValidType         string             `json:"valid_type,optional"`
	ValidStart        *common.LocalTime  `json:"valid_start,optional"`
	ValidEnd          *common.LocalTime  `json:"valid_end,optional"`
	ValidDays         int                `json:"valid_days,optional"`
	Stackable         *int8              `json:"stackable,optional"`
	UserIdentity      string             `json:"user_identity,optional"`
	Channels          []string           `json:"channels,optional"`
	Status            string             `json:"status,optional"`
	Remark            string             `json:"remark,optional"`
	Scopes            []CouponScopeItem  `json:"scopes,optional"`
}

// CouponUpdateBodyReq path+body
type CouponUpdateBodyReq struct {
	Id                uint64            `path:"id"`
	Name              string            `json:"name"`
	CouponType        string            `json:"coupon_type"`
	ThresholdAmount   float64           `json:"threshold_amount,optional"`
	DiscountAmount    float64           `json:"discount_amount,optional"`
	DiscountRate      float64           `json:"discount_rate,optional"`
	MaxDiscountAmount float64           `json:"max_discount_amount,optional"`
	ScopeType         string            `json:"scope_type,optional"`
	TotalCount        int               `json:"total_count,optional"`
	PerUserLimit      int               `json:"per_user_limit,optional"`
	ValidType         string            `json:"valid_type,optional"`
	ValidStart        *common.LocalTime `json:"valid_start,optional"`
	ValidEnd          *common.LocalTime `json:"valid_end,optional"`
	ValidDays         int               `json:"valid_days,optional"`
	Stackable         *int8             `json:"stackable,optional"`
	UserIdentity      string            `json:"user_identity,optional"`
	Channels          []string          `json:"channels,optional"`
	Status            string            `json:"status,optional"`
	Remark            string            `json:"remark,optional"`
	Scopes            []CouponScopeItem `json:"scopes,optional"`
}

func (r *CouponUpdateBodyReq) ToCouponSaveReq() CouponSaveReq {
	return CouponSaveReq{
		Name: r.Name, CouponType: r.CouponType, ThresholdAmount: r.ThresholdAmount,
		DiscountAmount: r.DiscountAmount, DiscountRate: r.DiscountRate,
		MaxDiscountAmount: r.MaxDiscountAmount, ScopeType: r.ScopeType,
		TotalCount: r.TotalCount, PerUserLimit: r.PerUserLimit, ValidType: r.ValidType,
		ValidStart: r.ValidStart, ValidEnd: r.ValidEnd, ValidDays: r.ValidDays,
		Stackable: r.Stackable, UserIdentity: r.UserIdentity, Channels: r.Channels,
		Status: r.Status, Remark: r.Remark, Scopes: r.Scopes,
	}
}

// GrantCouponReq 发放优惠券
type GrantCouponReq struct {
	CouponID uint64   `json:"coupon_id"`
	UserIDs  []uint64 `json:"user_ids"`
}

// ClaimCouponBodyReq 用户领券
type ClaimCouponBodyReq struct {
	Id     uint64 `path:"id"`
	Source string `json:"source,optional"`
}

// MatchItem 券匹配商品行
type MatchItem struct {
	ProductID      uint64  `json:"product_id"`
	CategoryID     uint64  `json:"category_id,optional"`
	Amount         float64 `json:"amount"`
	SeckillEntryID uint64  `json:"seckill_entry_id,optional"`
}

// MatchCouponsReq 内部匹配可用券
type MatchCouponsReq struct {
	UserID       uint64      `json:"user_id"`
	ShopID       uint64      `json:"shop_id"`
	Items        []MatchItem `json:"items"`
	UserCouponID uint64      `json:"user_coupon_id,optional"`
}

// LockCouponReq 锁定用户券
type LockCouponReq struct {
	UserCouponID   uint64  `json:"user_coupon_id"`
	UserID         uint64  `json:"user_id"`
	OrderID        uint64  `json:"order_id"`
	DiscountAmount float64 `json:"discount_amount"`
}

// UnlockCouponReq 解锁用户券
type UnlockCouponReq struct {
	UserCouponID uint64 `json:"user_coupon_id"`
	OrderID      uint64 `json:"order_id"`
}

// RedeemCouponReq 核销用户券
type RedeemCouponReq struct {
	UserCouponID   uint64  `json:"user_coupon_id"`
	OrderID        uint64  `json:"order_id"`
	DiscountAmount float64 `json:"discount_amount"`
}

// ReturnCouponReq 退还用户券
type ReturnCouponReq struct {
	UserCouponID uint64 `json:"user_coupon_id"`
	OrderID      uint64 `json:"order_id"`
}

// OrderGiftCouponReq 下单赠券
type OrderGiftCouponReq struct {
	UserID uint64 `json:"user_id"`
	ShopID uint64 `json:"shop_id"`
}

type SlotTypeQueryReq struct {
	SlotType string `form:"slot_type,optional"`
	City     string `form:"city,optional"`
}

type LocalShopsReq struct {
	Page     int     `form:"page,default=1"`
	PageSize int     `form:"page_size,default=10"`
	Lat      float64 `form:"lat,optional"`
	Lng      float64 `form:"lng,optional"`
	Keyword  string  `form:"keyword,optional"`
	Sort     string  `form:"sort,optional"`
}

type MapGeocoderReq struct {
	Lat float64 `form:"lat"`
	Lng float64 `form:"lng"`
}

type ShopIdQueryReq struct {
	ShopId uint64 `form:"shop_id,optional"`
}

// ---- list / query (go-zero form tags) ----

type IdPageReq struct {
	Id       uint64 `path:"id"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type StatusPageReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
}

type CouponListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
	Keyword  string `form:"keyword,optional"`
}

type SlotTypePageReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	SlotType string `form:"slot_type,optional"`
}

type SlotOrderListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	ShopId   uint64 `form:"shop_id,optional"`
	SlotType string `form:"slot_type,optional"`
	Status   string `form:"status,optional"`
}

type ThemePackageListReq struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	ThemeSlotId uint64 `form:"theme_slot_id,optional"`
}

type ThemeOrderListReq struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	ShopId      uint64 `form:"shop_id,optional"`
	ThemeSlotId uint64 `form:"theme_slot_id,optional"`
}

type ShopListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
	Name     string `form:"name,optional"`
}

type SeckillEntryListReq struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	SessionId uint64 `form:"session_id,optional"`
}

type ListData = ListResp
