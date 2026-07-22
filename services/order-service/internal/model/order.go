package model

import "mymall/common"

const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusFailed    = "failed"
	OrderStatusCancelled = "cancelled"
	OrderStatusShipped   = "shipped"
	OrderStatusCompleted = "completed"
	OrderStatusReviewed  = "reviewed"
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
	ID              uint64            `db:"id" json:"id"`
	OrderNo         string            `db:"order_no" json:"order_no"`
	UserID          uint64            `db:"user_id" json:"user_id"`
	ShopID          uint64            `db:"shop_id" json:"shop_id"`
	TotalAmount     float64           `db:"total_amount" json:"total_amount"`
	GoodsAmount     float64           `db:"goods_amount" json:"goods_amount"`
	DiscountAmount  float64           `db:"discount_amount" json:"discount_amount"`
	PayAmount       float64           `db:"pay_amount" json:"pay_amount"`
	UserCouponID    uint64            `db:"user_coupon_id" json:"user_coupon_id"`
	ReceiverName    string            `db:"receiver_name" json:"receiver_name"`
	ReceiverPhone   string            `db:"receiver_phone" json:"receiver_phone"`
	ReceiverAddress string            `db:"receiver_address" json:"receiver_address"`
	ShipCompany     string            `db:"ship_company" json:"ship_company"`
	ShipNo          string            `db:"ship_no" json:"ship_no"`
	ShippedAt       *common.LocalTime `db:"shipped_at" json:"shipped_at,omitempty"`
	CompletedAt     *common.LocalTime `db:"completed_at" json:"completed_at,omitempty"`
	ReviewedAt      *common.LocalTime `db:"reviewed_at" json:"reviewed_at,omitempty"`
	Remark          string            `db:"remark" json:"remark"`
	Status          string            `db:"status" json:"status"`
	CreatedAt       common.LocalTime  `db:"created_at" json:"created_at"`
	UpdatedAt       common.LocalTime  `db:"updated_at" json:"updated_at"`
	Items           []OrderItem       `db:"-" json:"items,omitempty"`
	UserName        string            `db:"-" json:"user_name,omitempty"`
	ShopName        string            `db:"-" json:"shop_name,omitempty"`
}

func (Order) TableName() string { return "orders" }

type OrderItem struct {
	ID             uint64           `db:"id" json:"id"`
	OrderID        uint64           `db:"order_id" json:"order_id"`
	ProductID      uint64           `db:"product_id" json:"product_id"`
	SkuID          uint64           `db:"sku_id" json:"sku_id"`
	ProductName    string           `db:"product_name" json:"product_name"`
	SkuSnapshot    string           `db:"sku_snapshot" json:"sku_snapshot,omitempty"`
	Price          float64          `db:"price" json:"price"`
	Quantity       int              `db:"quantity" json:"quantity"`
	SeckillEntryID uint64           `db:"seckill_entry_id" json:"seckill_entry_id"`
	CreatedAt      common.LocalTime `db:"created_at" json:"created_at"`
}

func (OrderItem) TableName() string { return "order_items" }

type OrderAfterSale struct {
	ID          uint64           `db:"id" json:"id"`
	OrderID     uint64           `db:"order_id" json:"order_id"`
	OrderNo     string           `db:"order_no" json:"order_no"`
	UserID      uint64           `db:"user_id" json:"user_id"`
	ShopID      uint64           `db:"shop_id" json:"shop_id"`
	Type        string           `db:"type" json:"type"`
	Reason      string           `db:"reason" json:"reason"`
	Amount      float64          `db:"amount" json:"amount"`
	Status      string           `db:"status" json:"status"`
	AdminRemark string           `db:"admin_remark" json:"admin_remark"`
	HandledBy   uint64           `db:"handled_by" json:"handled_by"`
	CreatedAt   common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt   common.LocalTime `db:"updated_at" json:"updated_at"`
	UserName    string           `db:"-" json:"user_name,omitempty"`
	ShopName    string           `db:"-" json:"shop_name,omitempty"`
}

func (OrderAfterSale) TableName() string { return "order_after_sales" }
