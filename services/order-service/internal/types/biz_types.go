package types

// Hand-maintained API/biz DTOs (kept in sync with api/order.api).

type CreateOrderItem struct {
	ProductID      uint64 `json:"product_id"`
	SkuID          uint64 `json:"sku_id"`
	SkuSnapshot    string `json:"sku_snapshot,optional"`
	Quantity       int    `json:"quantity"`
	SeckillEntryID uint64 `json:"seckill_entry_id,optional"`
}

type CreateOrderReq struct {
	AddressID    uint64            `json:"address_id"`
	Items        []CreateOrderItem `json:"items"`
	UserCouponID uint64            `json:"user_coupon_id,optional"`
}

type CouponPreviewReq struct {
	Items        []CreateOrderItem `json:"items"`
	UserCouponID uint64            `json:"user_coupon_id,optional"`
}

type OrderStatusCountsResp struct {
	Counts map[string]int64 `json:"counts"`
}

type ShipBodyReq struct {
	Id          uint64 `path:"id"`
	ShipCompany string `json:"ship_company"`
	ShipNo      string `json:"ship_no"`
}

type RemarkBodyReq struct {
	Id     uint64 `path:"id"`
	Remark string `json:"remark"`
}

type CreateAfterSaleBodyReq struct {
	Id     uint64  `path:"id"`
	Type   string  `json:"type"`
	Reason string  `json:"reason"`
	Amount float64 `json:"amount"`
}

type HandleAfterSaleBodyReq struct {
	Id          uint64 `path:"id"`
	Action      string `json:"action"`
	AdminRemark string `json:"admin_remark,optional"`
}

type LogisticsSaveBodyReq struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Sort   int    `json:"sort,optional"`
	Status int    `json:"status,optional"`
}

// Legacy aliases used by app/shared until fully sunk.
type ShipReq struct {
	ShipCompany string `json:"ship_company"`
	ShipNo      string `json:"ship_no"`
}

type RemarkReq struct {
	Remark string `json:"remark"`
}

type CreateAfterSaleReq struct {
	Type   string  `json:"type"`
	Reason string  `json:"reason"`
	Amount float64 `json:"amount"`
}

type HandleAfterSaleReq struct {
	Action      string `json:"action"`
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

type CreateReviewBodyReq struct {
	Id          uint64   `path:"id"`
	Rating      int8     `json:"rating"`
	Content     string   `json:"content,optional"`
	IsAnonymous bool     `json:"is_anonymous,optional"`
	OrderItemID uint64   `json:"order_item_id,optional"`
	Images      []string `json:"images,optional"`
}

type ReviewReplyBodyReq struct {
	Id    uint64 `path:"id"`
	Reply string `json:"reply"`
}

type LogisticsUpdateBodyReq struct {
	Id     uint64 `path:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Sort   int    `json:"sort,optional"`
	Status int    `json:"status,optional"`
}

type LogisticsStatusBodyReq struct {
	Id     uint64 `path:"id"`
	Status int8   `json:"status"`
}
