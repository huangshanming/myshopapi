package model

import "mymall/common"

const (
	ThemeLinkNone     = "none"
	ThemeLinkShop     = "shop"
	ThemeLinkCategory = "category"
	ThemeLinkProduct  = "product"

	ThemeSlotOn  = "on"
	ThemeSlotOff = "off"

	ThemeOrderActive    = "active"
	ThemeOrderExpired   = "expired"
	ThemeOrderCancelled = "cancelled"

	ThemePayWallet = "wallet"
	ThemePayAdmin  = "admin"

	WalletLogThemeSlot = "theme_slot"
)

type HomepageThemeSlot struct {
	ID               uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SlotKey          string           `gorm:"column:slot_key;type:varchar(32);uniqueIndex;not null" json:"slot_key"`
	Position         int              `gorm:"column:position;uniqueIndex;not null" json:"position"`
	Name             string           `gorm:"column:name;type:varchar(100)" json:"name"`
	Desc             string           `gorm:"column:desc;type:varchar(255)" json:"desc"`
	CoverURL         string           `gorm:"column:cover_url;type:varchar(500)" json:"cover_url"`
	DefaultLinkType  string           `gorm:"column:default_link_type;type:varchar(16);default:none" json:"default_link_type"`
	DefaultLinkID    uint64           `gorm:"column:default_link_id;default:0" json:"default_link_id"`
	Status           string           `gorm:"column:status;type:varchar(16);default:on" json:"status"`
	Sort             int              `gorm:"column:sort" json:"sort"`
	CreatedAt        common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
	// 展示
	OccupiedUntil string `gorm:"-" json:"occupied_until,omitempty"`
	HasActive     bool   `gorm:"-" json:"has_active,omitempty"`
}

func (HomepageThemeSlot) TableName() string { return "homepage_theme_slots" }

type HomepageThemePackage struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ThemeSlotID  uint64           `gorm:"column:theme_slot_id;default:0;index" json:"theme_slot_id"`
	Name         string           `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Price        float64          `gorm:"column:price;type:decimal(12,2)" json:"price"`
	DurationDays int              `gorm:"column:duration_days" json:"duration_days"`
	Status       string           `gorm:"column:status;type:varchar(16);default:on" json:"status"`
	Sort         int              `gorm:"column:sort" json:"sort"`
	Remark       string           `gorm:"column:remark;type:varchar(255)" json:"remark"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (HomepageThemePackage) TableName() string { return "homepage_theme_packages" }

type HomepageThemeOrder struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID       uint64           `gorm:"column:shop_id;not null;index" json:"shop_id"`
	ThemeSlotID  uint64           `gorm:"column:theme_slot_id;not null;index" json:"theme_slot_id"`
	PackageID    uint64           `gorm:"column:package_id" json:"package_id"`
	Title        string           `gorm:"column:title;type:varchar(100)" json:"title"`
	Subtitle     string           `gorm:"column:subtitle;type:varchar(255)" json:"subtitle"`
	CoverURL     string           `gorm:"column:cover_url;type:varchar(500)" json:"cover_url"`
	LinkType     string           `gorm:"column:link_type;type:varchar(16)" json:"link_type"`
	LinkID       uint64           `gorm:"column:link_id" json:"link_id"`
	Amount       float64          `gorm:"column:amount;type:decimal(12,2)" json:"amount"`
	DurationDays int              `gorm:"column:duration_days" json:"duration_days"`
	StartAt      common.LocalTime `gorm:"column:start_at" json:"start_at"`
	EndAt        common.LocalTime `gorm:"column:end_at" json:"end_at"`
	Status       string           `gorm:"column:status;type:varchar(16);default:active" json:"status"`
	PaySource    string           `gorm:"column:pay_source;type:varchar(16)" json:"pay_source"`
	WalletLogID  uint64           `gorm:"column:wallet_log_id" json:"wallet_log_id"`
	OperatorID   uint64           `gorm:"column:operator_id" json:"operator_id"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
	// enrichment
	ShopName      string `gorm:"-" json:"shop_name,omitempty"`
	ThemeSlotName string `gorm:"-" json:"theme_slot_name,omitempty"`
	PackageName   string `gorm:"-" json:"package_name,omitempty"`
}

func (HomepageThemeOrder) TableName() string { return "homepage_theme_orders" }

// ThemeTile C 端拼装卡片
type ThemeTile struct {
	Position  int    `json:"position"`
	SlotID    uint64 `json:"slot_id"`
	SlotKey   string `json:"slot_key"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	CoverURL  string `json:"cover_url"`
	LinkType  string `json:"link_type"`
	LinkID    uint64 `json:"link_id"`
	Paid      bool   `json:"paid"`
	ShopID    uint64 `json:"shop_id,omitempty"`
}
