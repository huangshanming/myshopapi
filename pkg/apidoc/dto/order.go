package dto

// CreateOrderItem 下单商品项
type CreateOrderItem struct {
	ProductID uint64 `json:"product_id" binding:"required" example:"1"`
	Quantity  int    `json:"quantity" binding:"required,min=1" example:"2"`
}

// CreateOrderReq 创建订单请求
type CreateOrderReq struct {
	Items []CreateOrderItem `json:"items" binding:"required,min=1,dive"`
}

// OrderItem 订单明细
type OrderItem struct {
	ID          uint64  `json:"id" example:"1"`
	OrderID     uint64  `json:"order_id" example:"1"`
	ProductID   uint64  `json:"product_id" example:"1"`
	ProductName string  `json:"product_name" example:"猫粮 10kg"`
	Price       float64 `json:"price" example:"159.00"`
	Quantity    int     `json:"quantity" example:"2"`
	CreatedAt   string  `json:"created_at" example:"2026-01-01 12:00:00"`
}

// OrderInfo 订单信息
type OrderInfo struct {
	ID          uint64      `json:"id" example:"1"`
	OrderNo     string      `json:"order_no" example:"ORDabc12345"`
	UserID      uint64      `json:"user_id" example:"1"`
	TotalAmount float64     `json:"total_amount" example:"318.00"`
	Status      string      `json:"status" example:"pending" enums:"pending,confirmed,failed,cancelled"`
	Items       []OrderItem `json:"items,omitempty"`
	CreatedAt   string      `json:"created_at" example:"2026-01-01 12:00:00"`
	UpdatedAt   string      `json:"updated_at" example:"2026-01-01 12:00:00"`
}

// OrderListResp 订单列表 data
type OrderListResp struct {
	Total int64       `json:"total" example:"10"`
	List  []OrderInfo `json:"list"`
}
