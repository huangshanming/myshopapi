package dto

// ProductListItem 商品列表项
type ProductListItem struct {
	ID          uint64  `json:"id" example:"1"`
	ProductNo   string  `json:"product_no" example:"P20260001"`
	Name        string  `json:"name" example:"猫粮 10kg"`
	Subtitle    string  `json:"subtitle,omitempty" example:"全年龄段"`
	MainImage   string  `json:"main_image,omitempty" example:"https://example.com/img.jpg"`
	MarketPrice float64 `json:"market_price,omitempty" example:"199.00"`
	SalePrice   float64 `json:"sale_price" example:"159.00"`
	Discount    float64 `json:"discount" example:"100.00"`
	Stock       int     `json:"stock" example:"100"`
	SoldCount   int     `json:"sold_count" example:"10"`
	CategoryID  uint64  `json:"category_id" example:"1"`
	PetType     string  `json:"pet_type" example:"cat"`
	Status      string  `json:"status" example:"on_sale"`
	IsHot       bool    `json:"is_hot" example:"true"`
	IsNew       bool    `json:"is_new" example:"false"`
	IsRecommend bool    `json:"is_recommend" example:"true"`
	CreatedAt   string  `json:"created_at" example:"2026-01-01 12:00:00"`
}

// ProductListResp 商品列表响应 data
type ProductListResp struct {
	Total int64             `json:"total" example:"100"`
	Data  []ProductListItem `json:"data"`
}

// ProductDetail 商品详情
type ProductDetail struct {
	ID          uint64  `json:"id" example:"1"`
	ProductNo   string  `json:"product_no" example:"P20260001"`
	Name        string  `json:"name" example:"猫粮 10kg"`
	Subtitle    string  `json:"subtitle,omitempty"`
	Description string  `json:"description,omitempty"`
	MainImage   string  `json:"main_image,omitempty"`
	SalePrice   float64 `json:"sale_price" example:"159.00"`
	MarketPrice float64 `json:"market_price,omitempty"`
	Stock       int     `json:"stock" example:"100"`
	CategoryID  uint64  `json:"category_id" example:"1"`
	PetType     string  `json:"pet_type" example:"cat"`
	Status      string  `json:"status" example:"on_sale"`
	CreatedAt   string  `json:"created_at" example:"2026-01-01 12:00:00"`
}

// CategoryItem 分类项
type CategoryItem struct {
	ID           uint64 `json:"id" example:"1"`
	ParentID     uint64 `json:"parent_id" example:"0"`
	Name         string `json:"name" example:"猫粮"`
	Icon         string `json:"icon,omitempty"`
	Image        string `json:"image,omitempty"`
	Description  string `json:"description,omitempty"`
	SortOrder    int    `json:"sort_order" example:"0"`
	Level        int    `json:"level" example:"1"`
	IsShow       bool   `json:"is_show" example:"true"`
	ProductCount int    `json:"product_count" example:"10"`
	CreatedAt    string `json:"created_at" example:"2026-01-01 12:00:00"`
}

// CategoryListResp 分类分页列表 data
type CategoryListResp struct {
	Total     int64          `json:"total" example:"20"`
	Page      int            `json:"page" example:"1"`
	PageSize  int            `json:"page_size" example:"10"`
	TotalPage int            `json:"total_page" example:"2"`
	List      []CategoryItem `json:"list"`
}

// CategoryDetailResp 分类详情
type CategoryDetailResp struct {
	CategoryItem
}
