package types

// CreateOrderItem 下单商品项
type CreateOrderItem struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
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
