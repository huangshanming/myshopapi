package types

// CreateOrderItem 下单商品项
type CreateOrderItem struct {
	ProductID   uint64 `json:"product_id"`
	SkuID       uint64 `json:"sku_id"`
	SkuSnapshot string `json:"sku_snapshot"` // 可选；不传则由服务端用默认 SKU 规格填充
	Quantity    int    `json:"quantity"`
}

// CreateOrderReq 创建订单请求
type CreateOrderReq struct {
	Items []CreateOrderItem `json:"items"`
}

// PageListResp 分页列表
type PageListResp struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
