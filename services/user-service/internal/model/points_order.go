package model

import "mymall/common"

const (
	PointsOrderPending   = "pending"
	PointsOrderShipped   = "shipped"
	PointsOrderCompleted = "completed"
	PointsOrderCancelled = "cancelled"

	PointsOrderRefType = "points_order"
)

type PointsExchangeOrder struct {
	ID              uint64            `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	OrderNo         string            `gorm:"column:order_no;type:varchar(32);uniqueIndex" db:"order_no" json:"order_no"`
	UserID          uint64            `gorm:"column:user_id;index" db:"user_id" json:"user_id"`
	ProductID       uint64            `gorm:"column:product_id" db:"product_id" json:"product_id"`
	ProductName     string            `gorm:"column:product_name;type:varchar(100)" db:"product_name" json:"product_name"`
	ProductCover    string            `gorm:"column:product_cover;type:varchar(512)" db:"product_cover" json:"product_cover"`
	Quantity        int               `gorm:"column:quantity" db:"quantity" json:"quantity"`
	PointsCost      int               `gorm:"column:points_cost" db:"points_cost" json:"points_cost"`
	Status          string            `gorm:"column:status;type:varchar(16)" db:"status" json:"status"`
	ReceiverName    string            `gorm:"column:receiver_name;type:varchar(64)" db:"receiver_name" json:"receiver_name"`
	ReceiverPhone   string            `gorm:"column:receiver_phone;type:varchar(32)" db:"receiver_phone" json:"receiver_phone"`
	ReceiverAddress string            `gorm:"column:receiver_address;type:varchar(255)" db:"receiver_address" json:"receiver_address"`
	ShipCompany     string            `gorm:"column:ship_company;type:varchar(64)" db:"ship_company" json:"ship_company"`
	ShipNo          string            `gorm:"column:ship_no;type:varchar(64)" db:"ship_no" json:"ship_no"`
	AdminRemark     string            `gorm:"column:admin_remark;type:varchar(255)" db:"admin_remark" json:"admin_remark"`
	ShippedAt       *common.LocalTime `gorm:"column:shipped_at" db:"shipped_at" json:"shipped_at,omitempty"`
	CompletedAt     *common.LocalTime `gorm:"column:completed_at" db:"completed_at" json:"completed_at,omitempty"`
	CancelledAt     *common.LocalTime `gorm:"column:cancelled_at" db:"cancelled_at" json:"cancelled_at,omitempty"`
	CreatedAt       common.LocalTime  `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt       common.LocalTime  `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (PointsExchangeOrder) TableName() string { return "points_exchange_orders" }
