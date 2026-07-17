package types

// SpecItem 规格项（如 颜色=[红,蓝]）
type SpecItem struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// SkuInput 前端提交的 SKU 行
type SkuInput struct {
	ID          uint64            `json:"id"`
	SpecValues  map[string]string `json:"spec_values"`
	SalePrice   float64           `json:"sale_price"`
	MarketPrice float64           `json:"market_price"`
	CostPrice   float64           `json:"cost_price"`
	Stock       int               `json:"stock"`
	StockWarn   int               `json:"stock_warn"`
	Barcode     string            `json:"barcode"`
	Status      string            `json:"status"`
}

// ImageInput 商品图片
type ImageInput struct {
	URL  string `json:"url"`
	Typ  string `json:"typ"` // main|gallery|detail
	Sort int    `json:"sort"`
}

// AttrInput 商品参数
type AttrInput struct {
	TemplateID uint64 `json:"template_id"`
	AttrKey    string `json:"attr_key"`
	AttrLabel  string `json:"attr_label"`
	AttrValue  string `json:"attr_value"`
}

// MerchantProductSaveReq 创建/更新商品（完整中台）
type MerchantProductSaveReq struct {
	Name        string       `json:"name"`
	Subtitle    string       `json:"subtitle"`
	Description string       `json:"description"`
	CategoryID  uint64       `json:"category_id"`
	ProductType string       `json:"product_type"` // physical|fresh|virtual
	PetType     string       `json:"pet_type"`
	Status      string       `json:"status"` // draft|on_sale
	MainImage   string       `json:"main_image"`
	SpecJSON    []SpecItem   `json:"spec_json"`
	Skus        []SkuInput   `json:"skus"`
	Images      []ImageInput `json:"images"`
	Attrs       []AttrInput  `json:"attrs"`
	TagIDs      []uint64     `json:"tag_ids"`
	// 生鲜扩展
	ShelfLife        int64  `json:"shelf_life"`
	StorageCondition string `json:"storage_condition"`
	// 虚拟扩展（复用 description 即可）
	SalePrice float64 `json:"sale_price"` // 无规格时的默认价
	Stock     int     `json:"stock"`      // 无规格时的默认库存
}

// MerchantProductReq 兼容旧接口
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

type SetStatusReq struct {
	Status string `json:"status"`
}

type CategoryReq struct {
	ParentId    uint64 `json:"parent_id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Level       int    `json:"level"`
	IsShow      *bool  `json:"is_show"`
}

type BatchProductReq struct {
	Action     string   `json:"action"` // on_sale|off_sale|category|price|recycle
	ProductIDs []uint64 `json:"product_ids"`
	CategoryID uint64   `json:"category_id"`
	PriceDelta float64  `json:"price_delta"` // 调价：加减金额
	PriceRate  float64  `json:"price_rate"`  // 或按比例，如 0.9
}

type StockAdjustReq struct {
	SkuID uint64 `json:"sku_id"`
	Stock *int   `json:"stock"` // 绝对库存
	Delta *int   `json:"delta"` // 相对增减
}

type BatchStockReq struct {
	Items []StockAdjustReq `json:"items"`
}

type RecycleReq struct {
	ProductIDs []uint64 `json:"product_ids"`
}

type ScheduleReq struct {
	Action string `json:"action"` // on_sale|off_sale
	RunAt  string `json:"run_at"` // 2006-01-02 15:04:05
}

type TagReq struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type AttrTemplateReq struct {
	Name      string `json:"name"`
	AttrsJSON string `json:"attrs_json"`
}
