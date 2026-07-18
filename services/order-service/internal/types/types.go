package types

// CreateOrderItem 下单商品项
type CreateOrderItem struct {
	ProductID   uint64 `json:"product_id"`
	SkuID       uint64 `json:"sku_id"`
	SkuSnapshot string `json:"sku_snapshot"`
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

type ShipReq struct {
	ShipCompany string `json:"ship_company"`
	ShipNo      string `json:"ship_no"`
}

type RemarkReq struct {
	Remark string `json:"remark"`
}

type CreateAfterSaleReq struct {
	Type   string  `json:"type"` // refund | return_refund
	Reason string  `json:"reason"`
	Amount float64 `json:"amount"`
}

type HandleAfterSaleReq struct {
	Action      string `json:"action"` // approve | reject | refunded | closed
	AdminRemark string `json:"admin_remark"`
}

type LogisticsSaveReq struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Sort   int    `json:"sort"`
	Status *int8  `json:"status"`
}

type LogisticsStatusReq struct {
	Status int8 `json:"status"`
}
