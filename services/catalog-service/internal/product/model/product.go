package model

import "mymall/common"

type BaseModel struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt common.LocalTime `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

type Product struct {
	BaseModel
	ShopID      uint64  `gorm:"column:shop_id;type:bigint;not null;default:0;index" json:"shop_id"`
	ProductNo   string  `gorm:"column:product_no;type:varchar(50);not null;uniqueIndex" json:"product_no"`
	Name        string  `gorm:"column:name;type:varchar(200);not null" json:"name"`
	Subtitle    string  `gorm:"column:subtitle;type:varchar(500);default:null" json:"subtitle,omitempty"`
	Description string  `gorm:"column:description;type:text;default:null" json:"description,omitempty"`
	MainImage   string  `gorm:"column:main_image;type:varchar(500);default:null" json:"main_image,omitempty"`
	ImageList   []uint8 `gorm:"column:image_list;type:json;default:null" json:"image_list,omitempty"`
	VideoUrl    string  `gorm:"column:video_url;type:varchar(500);default:null" json:"video_url,omitempty"`

	MarketPrice   float64 `gorm:"column:market_price;type:decimal(10,2);default:null" json:"market_price,omitempty"`
	SalePrice     float64 `gorm:"column:sale_price;type:decimal(10,2);not null" json:"sale_price"`
	CostPrice     float64 `gorm:"column:cost_price;type:decimal(10,2);default:null" json:"cost_price,omitempty"`
	Discount      float64 `gorm:"column:discount;type:decimal(5,2);default:100.00" json:"discount"`
	DiscountPrice float64 `gorm:"column:discount_price;type:decimal(10,2);default:null" json:"discount_price,omitempty"`

	Stock        int `gorm:"column:stock;type:int;default:0" json:"stock"`
	StockWarn    int `gorm:"column:stock_warn;type:int;default:10" json:"stock_warn"`
	SoldCount    int `gorm:"column:sold_count;type:int;default:0" json:"sold_count"`
	ViewCount    int `gorm:"column:view_count;type:int;default:0" json:"view_count"`
	CollectCount int `gorm:"column:collect_count;type:int;default:0" json:"collect_count"`

	PetType string  `gorm:"column:pet_type;type:enum('dog','cat','both','other');default:both;index" json:"pet_type"`
	PetAge  []uint8 `gorm:"column:pet_age;type:json;default:null" json:"pet_age,omitempty"`
	PetSize []uint8 `gorm:"column:pet_size;type:json;default:null" json:"pet_size,omitempty"`

	Weight      float64 `gorm:"column:weight;type:decimal(8,2);default:null" json:"weight,omitempty"`
	Unit        string  `gorm:"column:unit;type:varchar(20);default:null" json:"unit,omitempty"`
	BrandID     uint64  `gorm:"column:brand_id;type:int;default:null;index" json:"brand_id,omitempty"`
	CategoryID  uint64  `gorm:"column:category_id;type:int;not null;index" json:"category_id"`
	ProductType string  `gorm:"column:product_type;type:enum('physical','fresh','virtual');default:physical" json:"product_type"`
	SpecJSON    string  `gorm:"column:spec_json;type:json" json:"spec_json,omitempty"`
	Tags        []uint8 `gorm:"column:tags;type:json;default:null" json:"tags,omitempty"`

	NutritionInfo    []uint8 `gorm:"column:nutrition_info;type:json;default:null" json:"nutrition_info,omitempty"`
	Ingredients      string  `gorm:"column:ingredients;type:text;default:null" json:"ingredients,omitempty"`
	FeedingGuide     string  `gorm:"column:feeding_guide;type:text;default:null" json:"feeding_guide,omitempty"`
	ShelfLife        int64   `gorm:"column:shelf_life;type:int;default:null" json:"shelf_life,omitempty"`
	StorageCondition string  `gorm:"column:storage_condition;type:varchar(200);default:null" json:"storage_condition,omitempty"`

	Status         string            `gorm:"column:status;type:enum('draft','on_sale','off_sale','deleted','pending','approved','rejected');default:draft;index" json:"status"`
	IsHot          bool              `gorm:"column:is_hot;type:tinyint(1);default:0;index" json:"is_hot"`
	IsNew          bool              `gorm:"column:is_new;type:tinyint(1);default:0;index" json:"is_new"`
	IsRecommend    bool              `gorm:"column:is_recommend;type:tinyint(1);default:0;index" json:"is_recommend"`
	IsPrescription bool              `gorm:"column:is_prescription;type:tinyint(1);default:0" json:"is_prescription"`
	IsImported     bool              `gorm:"column:is_imported;type:tinyint(1);default:0" json:"is_imported"`
	IsOrganic      bool              `gorm:"column:is_organic;type:tinyint(1);default:0" json:"is_organic"`
	IsGrainFree    bool              `gorm:"column:is_grain_free;type:tinyint(1);default:0" json:"is_grain_free"`
	PublishTime    common.LocalTime  `gorm:"column:publish_time;type:datetime;default:null" json:"publish_time,omitempty"`
	ScheduleOnAt   *common.LocalTime `gorm:"column:schedule_on_at" json:"schedule_on_at,omitempty"`
	ScheduleOffAt  *common.LocalTime `gorm:"column:schedule_off_at" json:"schedule_off_at,omitempty"`
	CopyFromID     *uint64           `gorm:"column:copy_from_id" json:"copy_from_id,omitempty"`
	DeletedAt      common.LocalTime  `gorm:"column:deleted_at;type:timestamp;default:null;index" json:"deleted_at"`
}

type ProductListResp struct {
	ID             uint64           `gorm:"column:id" json:"id"`
	ShopID         uint64           `gorm:"column:shop_id" json:"shop_id"`
	ProductNo      string           `gorm:"column:product_no" json:"product_no"`
	Name           string           `gorm:"column:name" json:"name"`
	Subtitle       string           `gorm:"column:subtitle" json:"subtitle,omitempty"`
	MainImage      string           `gorm:"column:main_image" json:"main_image,omitempty"`
	MarketPrice    float64          `gorm:"column:market_price" json:"market_price,omitempty"`
	SalePrice      float64          `gorm:"column:sale_price" json:"sale_price"`
	Discount       float64          `gorm:"column:discount" json:"discount"`
	Stock          int              `gorm:"column:stock" json:"stock"`
	SoldCount      int              `gorm:"column:sold_count" json:"sold_count"`
	CategoryID     uint64           `gorm:"column:category_id" json:"category_id"`
	BrandID        int64            `gorm:"column:brand_id" json:"brand_id,omitempty"`
	PetType        string           `gorm:"column:pet_type" json:"pet_type"`
	Tags           []uint8          `gorm:"column:tags;type:json" json:"tags,omitempty"`
	Status         string           `gorm:"column:status" json:"status"`
	IsHot          bool             `gorm:"column:is_hot" json:"is_hot"`
	IsNew          bool             `gorm:"column:is_new" json:"is_new"`
	IsRecommend    bool             `gorm:"column:is_recommend" json:"is_recommend"`
	IsPrescription bool             `gorm:"column:is_prescription" json:"is_prescription,omitempty"`
	PublishTime    common.LocalTime `gorm:"column:publish_time" json:"publish_time,omitempty"`
	CreatedAt      common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (Product) TableName() string {
	return "products"
}

type ProductCategory struct {
	BaseModel
	ParentId     uint64 `gorm:"column:parent_id;type:int;default:0" json:"parent_id"`
	Name         string `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Icon         string `gorm:"column:icon;type:varchar(200);default:null" json:"icon,omitempty"`
	Image        string `gorm:"column:image;type:varchar(500);default:null" json:"image,omitempty"`
	Description  string `gorm:"column:description;type:varchar(500);default:null" json:"description,omitempty"`
	SortOrder    int    `gorm:"column:sort_order;type:int;default:0" json:"sort_order"`
	Level        int    `gorm:"column:level;type:int;default:1" json:"level"`
	IsShow       bool   `gorm:"column:is_show;type:tinyint(1);default:1" json:"is_show"`
	ProductCount int    `gorm:"column:product_count;type:int;default:0" json:"product_count"`
}

func (ProductCategory) TableName() string {
	return "product_categories"
}
