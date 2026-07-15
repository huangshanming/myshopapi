package types

// MerchantProductReq 商家创建/更新商品
type MerchantProductReq struct {
	Name       string  `json:"name"`
	SalePrice  float64 `json:"sale_price"`
	Stock      int     `json:"stock"`
	CategoryID uint64  `json:"category_id"`
	Subtitle   string  `json:"subtitle"`
	MainImage  string  `json:"main_image"`
	Status     string  `json:"status"`
	PetType    string  `json:"pet_type"`
}

// SetStatusReq 商家设置商品状态
type SetStatusReq struct {
	Status string `json:"status"`
}

// CategoryReq 平台创建/更新分类
type CategoryReq struct {
	ParentId    uint64 `json:"parent_id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Level       int    `json:"level"`
	IsShow      *bool  `json:"is_show"`
}
