package model

import "mymall/common"

const (
	UserWalletLogAdminAdjust   = "admin_adjust"
	UserWalletLogOrderFreeze   = "order_freeze"
	UserWalletLogOrderUnfreeze = "order_unfreeze"
	UserWalletLogOrderSettle   = "order_settle"

	UserWalletFieldBalance = "balance"
	UserWalletFieldFrozen  = "frozen_balance"
)

type UserWallet struct {
	UserID        uint64           `gorm:"column:user_id;primaryKey" db:"user_id" json:"user_id"`
	Balance       float64          `gorm:"column:balance;type:decimal(12,2)" db:"balance" json:"balance"`
	FrozenBalance float64          `gorm:"column:frozen_balance;type:decimal(12,2)" db:"frozen_balance" json:"frozen_balance"`
	CreatedAt     common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt     common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (UserWallet) TableName() string { return "user_wallets" }

type UserWalletLog struct {
	ID             uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	UserID         uint64           `gorm:"column:user_id;index" db:"user_id" json:"user_id"`
	ChangeType     string           `gorm:"column:change_type;type:varchar(32)" db:"change_type" json:"change_type"`
	Amount         float64          `gorm:"column:amount;type:decimal(12,2)" db:"amount" json:"amount"`
	BalanceAfter   float64          `gorm:"column:balance_after;type:decimal(12,2)" db:"balance_after" json:"balance_after"`
	FrozenAfter    float64          `gorm:"column:frozen_after;type:decimal(12,2)" db:"frozen_after" json:"frozen_after"`
	Remark         string           `gorm:"column:remark;type:varchar(255)" db:"remark" json:"remark"`
	OperatorUserID *uint64          `gorm:"column:operator_user_id" db:"operator_user_id" json:"operator_user_id,omitempty"`
	RefType        string           `gorm:"column:ref_type;type:varchar(32)" db:"ref_type" json:"ref_type"`
	RefID          uint64           `gorm:"column:ref_id" db:"ref_id" json:"ref_id"`
	CreatedAt      common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
}

func (UserWalletLog) TableName() string { return "user_wallet_logs" }
