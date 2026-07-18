package model

import "mymall/common"

const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusFailed    = "failed"
	OrderStatusCancelled = "cancelled"
	OrderStatusShipped   = "shipped"
	OrderStatusCompleted = "completed"
)

const (
	AfterSaleTypeRefund       = "refund"
	AfterSaleTypeReturnRefund = "return_refund"

	AfterSalePending  = "pending"
	AfterSaleApproved = "approved"
	AfterSaleRejected = "rejected"
	AfterSaleRefunded = "refunded"
	AfterSaleClosed   = "closed"
)

type StockItem struct {
	ProductID uint64 `json:"product_id"`
	SkuID     uint64 `json:"sku_id"`
	Quantity  int    `json:"quantity"`
}

type Order struct {
	ID              uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo         string           `gorm:"column:order_no;type:varchar(64);not null;uniqueIndex" json:"order_no"`
	UserID          uint64           `gorm:"column:user_id;not null;index" json:"user_id"`
	ShopID          uint64           `gorm:"column:shop_id;not null;default:0;index" json:"shop_id"`
	TotalAmount     float64          `gorm:"column:total_amount;type:decimal(12,2);not null" json:"total_amount"`
	ReceiverName    string           `gorm:"column:receiver_name;type:varchar(64);not null;default:''" json:"receiver_name"`
	ReceiverPhone   string           `gorm:"column:receiver_phone;type:varchar(20);not null;default:''" json:"receiver_phone"`
	ReceiverAddress string           `gorm:"column:receiver_address;type:varchar(255);not null;default:''" json:"receiver_address"`
	ShipCompany     string           `gorm:"column:ship_company;type:varchar(64);not null;default:''" json:"ship_company"`
	ShipNo          string           `gorm:"column:ship_no;type:varchar(64);not null;default:''" json:"ship_no"`
	ShippedAt       *common.LocalTime `gorm:"column:shipped_at" json:"shipped_at,omitempty"`
	CompletedAt     *common.LocalTime `gorm:"column:completed_at" json:"completed_at,omitempty"`
	Remark          string           `gorm:"column:remark;type:varchar(255);not null;default:''" json:"remark"`
	Status          string           `gorm:"column:status;type:enum('pending','confirmed','failed','cancelled','shipped','completed');default:pending" json:"status"`
	CreatedAt       common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
	Items           []OrderItem      `gorm:"foreignKey:OrderID" json:"items,omitempty"`

	// enrichment (not persisted)
	UserName string `gorm:"-" json:"user_name,omitempty"`
	ShopName string `gorm:"-" json:"shop_name,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID          uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID     uint64           `gorm:"column:order_id;not null;index" json:"order_id"`
	ProductID   uint64           `gorm:"column:product_id;not null" json:"product_id"`
	SkuID       uint64           `gorm:"column:sku_id;not null;default:0" json:"sku_id"`
	ProductName string           `gorm:"column:product_name;type:varchar(200);not null" json:"product_name"`
	SkuSnapshot string           `gorm:"column:sku_snapshot;type:json" json:"sku_snapshot,omitempty"`
	Price       float64          `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Quantity    int              `gorm:"column:quantity;not null" json:"quantity"`
	CreatedAt   common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

type OrderAfterSale struct {
	ID          uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID     uint64           `gorm:"column:order_id;not null;index" json:"order_id"`
	OrderNo     string           `gorm:"column:order_no;type:varchar(64);not null" json:"order_no"`
	UserID      uint64           `gorm:"column:user_id;not null;index" json:"user_id"`
	ShopID      uint64           `gorm:"column:shop_id;not null;default:0;index" json:"shop_id"`
	Type        string           `gorm:"column:type;type:enum('refund','return_refund');default:refund" json:"type"`
	Reason      string           `gorm:"column:reason;type:varchar(500);not null;default:''" json:"reason"`
	Amount      float64          `gorm:"column:amount;type:decimal(12,2);not null" json:"amount"`
	Status      string           `gorm:"column:status;type:enum('pending','approved','rejected','refunded','closed');default:pending;index" json:"status"`
	AdminRemark string           `gorm:"column:admin_remark;type:varchar(500);not null;default:''" json:"admin_remark"`
	HandledBy   uint64           `gorm:"column:handled_by;not null;default:0" json:"handled_by"`
	CreatedAt   common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   common.LocalTime `gorm:"column:updated_at" json:"updated_at"`

	UserName string `gorm:"-" json:"user_name,omitempty"`
	ShopName string `gorm:"-" json:"shop_name,omitempty"`
}

func (OrderAfterSale) TableName() string {
	return "order_after_sales"
}
