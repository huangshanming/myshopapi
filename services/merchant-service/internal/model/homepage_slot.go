package model

import "mymall/common"

const (
	SlotBrandShop   = "brand_shop"
	SlotQualityShop = "quality_shop"
	SlotArticle     = "article"

	SlotPkgOn  = "on"
	SlotPkgOff = "off"

	SlotOrderActive    = "active"
	SlotOrderExpired   = "expired"
	SlotOrderCancelled = "cancelled"

	SlotPayWallet = "wallet"
	SlotPayAdmin  = "admin"

	WalletLogHomepageSlot = "homepage_slot"
)

type HomepageSlotPackage struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SlotType     string           `gorm:"column:slot_type;type:varchar(32);not null;index" json:"slot_type"`
	Name         string           `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Price        float64          `gorm:"column:price;type:decimal(12,2)" json:"price"`
	DurationDays int              `gorm:"column:duration_days" json:"duration_days"`
	Status       string           `gorm:"column:status;type:varchar(16);default:on" json:"status"`
	Sort         int              `gorm:"column:sort" json:"sort"`
	Remark       string           `gorm:"column:remark;type:varchar(255)" json:"remark"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (HomepageSlotPackage) TableName() string { return "homepage_slot_packages" }

type HomepageSlotSetting struct {
	SlotType  string           `gorm:"column:slot_type;primaryKey;type:varchar(32)" json:"slot_type"`
	HomeLimit int              `gorm:"column:home_limit" json:"home_limit"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (HomepageSlotSetting) TableName() string { return "homepage_slot_settings" }

type HomepageSlotOrder struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID       uint64           `gorm:"column:shop_id;not null;index" json:"shop_id"`
	SlotType     string           `gorm:"column:slot_type;type:varchar(32);not null;index" json:"slot_type"`
	PackageID    uint64           `gorm:"column:package_id" json:"package_id"`
	TargetID     uint64           `gorm:"column:target_id;index" json:"target_id"`
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

	ShopName    string `gorm:"-" json:"shop_name,omitempty"`
	PackageName string `gorm:"-" json:"package_name,omitempty"`
	TargetName  string `gorm:"-" json:"target_name,omitempty"`
	Paid        bool   `gorm:"-" json:"paid,omitempty"`
}

func (HomepageSlotOrder) TableName() string { return "homepage_slot_orders" }
