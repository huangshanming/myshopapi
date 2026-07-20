package model

import "mymall/common"

const (
	WalletLogAdminAdjust   = "admin_adjust"
	WalletLogSeckillApply  = "seckill_apply"
	WalletLogSeckillRenew  = "seckill_renew"

	WalletFieldBalance = "balance"
	WalletFieldDeposit = "deposit"
	WalletFieldFrozen  = "frozen_balance"

	SeckillSessionActive = "active"
	SeckillSessionEnded  = "ended"

	SeckillEntryActive = "active"
	SeckillRuleOn      = 1
)

type ShopWallet struct {
	ShopID        uint64           `gorm:"column:shop_id;primaryKey" json:"shop_id"`
	Balance       float64          `gorm:"column:balance;type:decimal(12,2)" json:"balance"`
	FrozenBalance float64          `gorm:"column:frozen_balance;type:decimal(12,2)" json:"frozen_balance"`
	Deposit       float64          `gorm:"column:deposit;type:decimal(12,2)" json:"deposit"`
	CreatedAt     common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ShopWallet) TableName() string { return "shop_wallets" }

type ShopWalletLog struct {
	ID             uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID         uint64           `gorm:"column:shop_id;index" json:"shop_id"`
	ChangeType     string           `gorm:"column:change_type;type:varchar(32)" json:"change_type"`
	Amount         float64          `gorm:"column:amount;type:decimal(12,2)" json:"amount"`
	BalanceAfter   float64          `gorm:"column:balance_after;type:decimal(12,2)" json:"balance_after"`
	FrozenAfter    float64          `gorm:"column:frozen_after;type:decimal(12,2)" json:"frozen_after"`
	DepositAfter   float64          `gorm:"column:deposit_after;type:decimal(12,2)" json:"deposit_after"`
	Remark         string           `gorm:"column:remark;type:varchar(255)" json:"remark"`
	OperatorUserID *uint64          `gorm:"column:operator_user_id" json:"operator_user_id,omitempty"`
	RefType        string           `gorm:"column:ref_type;type:varchar(32)" json:"ref_type"`
	RefID          uint64           `gorm:"column:ref_id" json:"ref_id"`
	CreatedAt      common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (ShopWalletLog) TableName() string { return "shop_wallet_logs" }

type SeckillRule struct {
	ID                 uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DurationHours      int              `gorm:"column:duration_hours" json:"duration_hours"`
	ApplyFee           float64          `gorm:"column:apply_fee;type:decimal(12,2)" json:"apply_fee"`
	MaxEntriesPerShop  int              `gorm:"column:max_entries_per_shop" json:"max_entries_per_shop"`
	Status             int8             `gorm:"column:status" json:"status"`
	CreatedAt          common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (SeckillRule) TableName() string { return "seckill_rules" }

type SeckillSession struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RuleID    uint64           `gorm:"column:rule_id" json:"rule_id"`
	StartAt   common.LocalTime `gorm:"column:start_at" json:"start_at"`
	EndAt     common.LocalTime `gorm:"column:end_at" json:"end_at"`
	Status    string           `gorm:"column:status;type:varchar(16)" json:"status"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (SeckillSession) TableName() string { return "seckill_sessions" }

type SeckillEntry struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SessionID    uint64           `gorm:"column:session_id;index" json:"session_id"`
	ShopID       uint64           `gorm:"column:shop_id;index" json:"shop_id"`
	ProductID    uint64           `gorm:"column:product_id" json:"product_id"`
	ProductName  string           `gorm:"column:product_name;type:varchar(200)" json:"product_name"`
	ProductImage string           `gorm:"column:product_image;type:varchar(500)" json:"product_image"`
	OriginPrice  float64          `gorm:"column:origin_price;type:decimal(12,2)" json:"origin_price"`
	SeckillPrice float64          `gorm:"column:seckill_price;type:decimal(12,2)" json:"seckill_price"`
	SeckillStock int              `gorm:"column:seckill_stock" json:"seckill_stock"`
	FeeAmount    float64          `gorm:"column:fee_amount;type:decimal(12,2)" json:"fee_amount"`
	Status       string           `gorm:"column:status;type:varchar(16)" json:"status"`
	AutoRenew    int8             `gorm:"column:auto_renew" json:"auto_renew"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (SeckillEntry) TableName() string { return "seckill_entries" }
