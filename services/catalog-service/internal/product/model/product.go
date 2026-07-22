package model

import "mymall/common"

type BaseModel struct {
	ID        uint64           `db:"id" json:"id"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `db:"updated_at" json:"updated_at"`
}

type Product struct {
	BaseModel
	ShopID      uint64  `db:"shop_id" json:"shop_id"`
	ProductNo   string  `db:"product_no" json:"product_no"`
	Name        string  `db:"name" json:"name"`
	Subtitle    string  `db:"subtitle" json:"subtitle,omitempty"`
	Description string  `db:"description" json:"description,omitempty"`
	MainImage   string  `db:"main_image" json:"main_image,omitempty"`
	ImageList   []uint8 `db:"image_list" json:"image_list,omitempty"`
	VideoUrl    string  `db:"video_url" json:"video_url,omitempty"`

	MarketPrice   float64 `db:"market_price" json:"market_price,omitempty"`
	SalePrice     float64 `db:"sale_price" json:"sale_price"`
	CostPrice     float64 `db:"cost_price" json:"cost_price,omitempty"`
	Discount      float64 `db:"discount" json:"discount"`
	DiscountPrice float64 `db:"discount_price" json:"discount_price,omitempty"`

	Stock        int `db:"stock" json:"stock"`
	StockWarn    int `db:"stock_warn" json:"stock_warn"`
	SoldCount    int     `db:"sold_count" json:"sold_count"`
	ViewCount    int     `db:"view_count" json:"view_count"`
	CollectCount int     `db:"collect_count" json:"collect_count"`
	AvgRating    float64 `db:"avg_rating" json:"avg_rating"`
	ReviewCount  int     `db:"review_count" json:"review_count"`
	GoodRate     float64 `db:"good_rate" json:"good_rate"`

	PetType string  `db:"pet_type" json:"pet_type"`
	PetAge  []uint8 `db:"pet_age" json:"pet_age,omitempty"`
	PetSize []uint8 `db:"pet_size" json:"pet_size,omitempty"`

	Weight      float64 `db:"weight" json:"weight,omitempty"`
	Unit        string  `db:"unit" json:"unit,omitempty"`
	BrandID     uint64  `db:"brand_id" json:"brand_id,omitempty"`
	CategoryID  uint64  `db:"category_id" json:"category_id"`
	ProductType string  `db:"product_type" json:"product_type"`
	SpecJSON    string  `db:"spec_json" json:"spec_json,omitempty"`
	Tags        []uint8 `db:"tags" json:"tags,omitempty"`

	NutritionInfo    []uint8 `db:"nutrition_info" json:"nutrition_info,omitempty"`
	Ingredients      string  `db:"ingredients" json:"ingredients,omitempty"`
	FeedingGuide     string  `db:"feeding_guide" json:"feeding_guide,omitempty"`
	ShelfLife        int64   `db:"shelf_life" json:"shelf_life,omitempty"`
	StorageCondition string  `db:"storage_condition" json:"storage_condition,omitempty"`

	Status         string            `db:"status" json:"status"`
	IsHot          bool              `db:"is_hot" json:"is_hot"`
	IsNew          bool              `db:"is_new" json:"is_new"`
	IsRecommend    bool              `db:"is_recommend" json:"is_recommend"`
	IsPrescription bool              `db:"is_prescription" json:"is_prescription"`
	IsImported     bool              `db:"is_imported" json:"is_imported"`
	IsOrganic      bool              `db:"is_organic" json:"is_organic"`
	IsGrainFree    bool              `db:"is_grain_free" json:"is_grain_free"`
	PublishTime    common.LocalTime  `db:"publish_time" json:"publish_time,omitempty"`
	ScheduleOnAt   *common.LocalTime `db:"schedule_on_at" json:"schedule_on_at,omitempty"`
	ScheduleOffAt  *common.LocalTime `db:"schedule_off_at" json:"schedule_off_at,omitempty"`
	// Non-pointer: go-zero sqlx pre-allocates *uint64 so NULL cannot scan; use IFNULL in SELECT.
	CopyFromID     uint64            `db:"copy_from_id" json:"copy_from_id,omitempty"`
	DeletedAt      common.LocalTime  `db:"deleted_at" json:"deleted_at"`
}

type ProductListResp struct {
	ID             uint64           `db:"id" json:"id"`
	ShopID         uint64           `db:"shop_id" json:"shop_id"`
	ProductNo      string           `db:"product_no" json:"product_no"`
	Name           string           `db:"name" json:"name"`
	Subtitle       string           `db:"subtitle" json:"subtitle,omitempty"`
	MainImage      string           `db:"main_image" json:"main_image,omitempty"`
	MarketPrice    float64          `db:"market_price" json:"market_price,omitempty"`
	SalePrice      float64          `db:"sale_price" json:"sale_price"`
	Discount       float64          `db:"discount" json:"discount"`
	Stock          int              `db:"stock" json:"stock"`
	SoldCount      int              `db:"sold_count" json:"sold_count"`
	CollectCount   int              `db:"collect_count" json:"collect_count"`
	AvgRating      float64          `db:"avg_rating" json:"avg_rating"`
	ReviewCount    int              `db:"review_count" json:"review_count"`
	GoodRate       float64          `db:"good_rate" json:"good_rate"`
	CategoryID     uint64           `db:"category_id" json:"category_id"`
	BrandID        int64            `db:"brand_id" json:"brand_id,omitempty"`
	PetType        string           `db:"pet_type" json:"pet_type"`
	Tags           []uint8          `db:"tags" json:"tags,omitempty"`
	Status         string           `db:"status" json:"status"`
	IsHot          bool             `db:"is_hot" json:"is_hot"`
	IsNew          bool             `db:"is_new" json:"is_new"`
	IsRecommend    bool             `db:"is_recommend" json:"is_recommend"`
	IsPrescription bool             `db:"is_prescription" json:"is_prescription,omitempty"`
	PublishTime    common.LocalTime `db:"publish_time" json:"publish_time,omitempty"`
	CreatedAt      common.LocalTime `db:"created_at" json:"created_at"`
}

func (Product) TableName() string {
	return "products"
}

// ProductSalesRankItem 销量榜：优先今日销量，再总销量
type ProductSalesRankItem struct {
	ID        uint64  `json:"id"`
	ShopID    uint64  `json:"shop_id"`
	Name      string  `json:"name"`
	ShopName  string  `json:"shop_name"`
	MainImage string  `json:"main_image"`
	SalePrice float64 `json:"sale_price"`
	SoldCount int     `json:"sold_count"`
	TodaySold int     `json:"today_sold"`
}

type ProductCategory struct {
	BaseModel
	ParentId     uint64 `db:"parent_id" json:"parent_id"`
	Name         string `db:"name" json:"name"`
	Icon         string `db:"icon" json:"icon,omitempty"`
	Image        string `db:"image" json:"image,omitempty"`
	Description  string `db:"description" json:"description,omitempty"`
	SortOrder    int    `db:"sort_order" json:"sort_order"`
	Level        int    `db:"level" json:"level"`
	IsShow       bool   `db:"is_show" json:"is_show"`
	ProductCount int    `db:"product_count" json:"product_count"`
}

func (ProductCategory) TableName() string {
	return "product_categories"
}
