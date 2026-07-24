package types

// Legacy / biz-only DTOs — must NOT duplicate names from api/order.api → types.go.
// After ./scripts/gen-api.sh, API request/response types live only in types.go.

// Legacy aliases used by biz until callers use *BodyReq directly.
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

// Aliases — canonical resp types live in types.go
type OrderDetailData = OrderDetailResp
type ListData = ListResp
