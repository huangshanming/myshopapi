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
	ID              uint64            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo         string            `gorm:"column:order_no;type:varchar(32);uniqueIndex" json:"order_no"`
	UserID          uint64            `gorm:"column:user_id;index" json:"user_id"`
	ProductID       uint64            `gorm:"column:product_id" json:"product_id"`
	ProductName     string            `gorm:"column:product_name;type:varchar(100)" json:"product_name"`
	ProductCover    string            `gorm:"column:product_cover;type:varchar(512)" json:"product_cover"`
	Quantity        int               `gorm:"column:quantity" json:"quantity"`
	PointsCost      int               `gorm:"column:points_cost" json:"points_cost"`
	Status          string            `gorm:"column:status;type:varchar(16)" json:"status"`
	ReceiverName    string            `gorm:"column:receiver_name;type:varchar(64)" json:"receiver_name"`
	ReceiverPhone   string            `gorm:"column:receiver_phone;type:varchar(32)" json:"receiver_phone"`
	ReceiverAddress string            `gorm:"column:receiver_address;type:varchar(255)" json:"receiver_address"`
	ShipCompany     string            `gorm:"column:ship_company;type:varchar(64)" json:"ship_company"`
	ShipNo          string            `gorm:"column:ship_no;type:varchar(64)" json:"ship_no"`
	AdminRemark     string            `gorm:"column:admin_remark;type:varchar(255)" json:"admin_remark"`
	ShippedAt       *common.LocalTime `gorm:"column:shipped_at" json:"shipped_at,omitempty"`
	CompletedAt     *common.LocalTime `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CancelledAt     *common.LocalTime `gorm:"column:cancelled_at" json:"cancelled_at,omitempty"`
	CreatedAt       common.LocalTime  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       common.LocalTime  `gorm:"column:updated_at" json:"updated_at"`
}

func (PointsExchangeOrder) TableName() string { return "points_exchange_orders" }
