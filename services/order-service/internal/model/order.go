package model

import "mymall/common"

const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusFailed    = "failed"
	OrderStatusCancelled = "cancelled"
)

type StockItem struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Order struct {
	ID          uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo     string           `gorm:"column:order_no;type:varchar(64);not null;uniqueIndex" json:"order_no"`
	UserID      uint64           `gorm:"column:user_id;not null;index" json:"user_id"`
	ShopID      uint64           `gorm:"column:shop_id;not null;default:0;index" json:"shop_id"`
	TotalAmount float64          `gorm:"column:total_amount;type:decimal(12,2);not null" json:"total_amount"`
	Status      string           `gorm:"column:status;type:enum('pending','confirmed','failed','cancelled');default:pending" json:"status"`
	CreatedAt   common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
	Items       []OrderItem      `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID          uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID     uint64           `gorm:"column:order_id;not null;index" json:"order_id"`
	ProductID   uint64           `gorm:"column:product_id;not null" json:"product_id"`
	ProductName string           `gorm:"column:product_name;type:varchar(200);not null" json:"product_name"`
	Price       float64          `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Quantity    int              `gorm:"column:quantity;not null" json:"quantity"`
	CreatedAt   common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
