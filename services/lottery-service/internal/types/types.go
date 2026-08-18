package types

type EmptyResp struct{}

type PageReq struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

type PageListResp struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

type DataResp struct {
	Data interface{} `json:"data,optional"`
}

type URLResp struct {
	URL string `json:"url"`
}

type ActivitySaveReq struct {
	Title      string `json:"title"`
	Status     int    `json:"status,optional"`
	CostPoints int    `json:"cost_points,optional"`
	DailyLimit int    `json:"daily_limit,optional"`
	StartAt    string `json:"start_at,optional"`
	EndAt      string `json:"end_at,optional"`
}

type ActivityUpdateReq struct {
	Id         uint64 `path:"id"`
	Title      string `json:"title"`
	Status     int    `json:"status,optional"`
	CostPoints int    `json:"cost_points,optional"`
	DailyLimit int    `json:"daily_limit,optional"`
	StartAt    string `json:"start_at,optional"`
	EndAt      string `json:"end_at,optional"`
}

type IdPathReq struct {
	Id uint64 `path:"id"`
}

type PrizeItemReq struct {
	Slot         int    `json:"slot"`
	Name         string `json:"name"`
	CoverURL     string `json:"cover_url,optional"`
	PrizeType    string `json:"prize_type"`
	PointsAmount int    `json:"points_amount,optional"`
	Weight       int    `json:"weight"`
	Stock        int    `json:"stock"`
	StockStrict  int    `json:"stock_strict,optional"`
}

type PrizesSaveReq struct {
	Id     uint64         `path:"id"`
	Prizes []PrizeItemReq `json:"prizes"`
}

type AdminRecordsReq struct {
	ActivityID uint64 `form:"activity_id,optional"`
	PrizeType  string `form:"prize_type,optional"` // points | thanks | physical；空=全部
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=20"`
}

type ClaimAddressReq struct {
	Id        uint64 `path:"id"`
	AddressID uint64 `json:"address_id"`
}

type AdminOrdersReq struct {
	FulfillStatus string `form:"fulfill_status,optional"`
	Page          int    `form:"page,default=1"`
	PageSize      int    `form:"page_size,default=20"`
}

type AdminShipReq struct {
	Id          uint64 `path:"id"`
	ShipCompany string `json:"ship_company,optional"`
	ShipNo      string `json:"ship_no"`
}
