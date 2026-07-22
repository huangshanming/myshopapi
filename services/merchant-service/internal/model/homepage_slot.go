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
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	SlotType     string           `gorm:"column:slot_type;type:varchar(32);not null;index" db:"slot_type" json:"slot_type"`
	Name         string           `gorm:"column:name;type:varchar(100);not null" db:"name" json:"name"`
	Price        float64          `gorm:"column:price;type:decimal(12,2)" db:"price" json:"price"`
	DurationDays int              `gorm:"column:duration_days" db:"duration_days" json:"duration_days"`
	Status       string           `gorm:"column:status;type:varchar(16);default:on" db:"status" json:"status"`
	Sort         int              `gorm:"column:sort" db:"sort" json:"sort"`
	Remark       string           `gorm:"column:remark;type:varchar(255)" db:"remark" json:"remark"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (HomepageSlotPackage) TableName() string { return "homepage_slot_packages" }

type HomepageSlotSetting struct {
	SlotType  string           `gorm:"column:slot_type;primaryKey;type:varchar(32)" db:"slot_type" json:"slot_type"`
	HomeLimit int              `gorm:"column:home_limit" db:"home_limit" json:"home_limit"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (HomepageSlotSetting) TableName() string { return "homepage_slot_settings" }

type HomepageSlotOrder struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	ShopID       uint64           `gorm:"column:shop_id;not null;index" db:"shop_id" json:"shop_id"`
	SlotType     string           `gorm:"column:slot_type;type:varchar(32);not null;index" db:"slot_type" json:"slot_type"`
	PackageID    uint64           `gorm:"column:package_id" db:"package_id" json:"package_id"`
	TargetID     uint64           `gorm:"column:target_id;index" db:"target_id" json:"target_id"`
	Amount       float64          `gorm:"column:amount;type:decimal(12,2)" db:"amount" json:"amount"`
	DurationDays int              `gorm:"column:duration_days" db:"duration_days" json:"duration_days"`
	StartAt      common.LocalTime `gorm:"column:start_at" db:"start_at" json:"start_at"`
	EndAt        common.LocalTime `gorm:"column:end_at" db:"end_at" json:"end_at"`
	Status       string           `gorm:"column:status;type:varchar(16);default:active" db:"status" json:"status"`
	PaySource    string           `gorm:"column:pay_source;type:varchar(16)" db:"pay_source" json:"pay_source"`
	WalletLogID  uint64           `gorm:"column:wallet_log_id" db:"wallet_log_id" json:"wallet_log_id"`
	OperatorID   uint64           `gorm:"column:operator_id" db:"operator_id" json:"operator_id"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`

	ShopName    string `gorm:"-" db:"-" json:"shop_name,omitempty"`
	PackageName string `gorm:"-" db:"-" json:"package_name,omitempty"`
	TargetName  string `gorm:"-" db:"-" json:"target_name,omitempty"`
	Paid        bool   `gorm:"-" db:"-" json:"paid,omitempty"`
}

func (HomepageSlotOrder) TableName() string { return "homepage_slot_orders" }
