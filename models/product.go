package model

import (
	"mymall/common"
)

// Product 商品主表模型，对应数据库表products
// 字段类型严格匹配数据库表定义，标签包含GORM配置、JSON序列化配置及注释
type Product struct {
	// 基础字段
	BaseModel
	ProductNo   string  `gorm:"column:product_no;type:varchar(50);not null;uniqueIndex;comment:商品编号" json:"product_no"`
	Name        string  `gorm:"column:name;type:varchar(200);not null;comment:商品名称" json:"name"`
	Subtitle    string  `gorm:"column:subtitle;type:varchar(500);default:null;comment:商品副标题" json:"subtitle,omitempty"` // 可空字段使用sql.NullString避免空指针
	Description string  `gorm:"column:description;type:text;default:null;comment:商品描述" json:"description,omitempty"`
	MainImage   string  `gorm:"column:main_image;type:varchar(500);default:null;comment:主图" json:"main_image,omitempty"`
	ImageList   []uint8 `gorm:"column:image_list;type:json;default:null;comment:商品图片列表" json:"image_list,omitempty"` // JSON数组映射为string切片
	VideoUrl    string  `gorm:"column:video_url;type:varchar(500);default:null;comment:商品视频" json:"video_url,omitempty"`

	// 价格相关字段
	MarketPrice   float64 `gorm:"column:market_price;type:decimal(10,2);default:null;comment:市场价" json:"market_price,omitempty"`
	SalePrice     float64 `gorm:"column:sale_price;type:decimal(10,2);not null;comment:销售价" json:"sale_price"`
	CostPrice     float64 `gorm:"column:cost_price;type:decimal(10,2);default:null;comment:成本价" json:"cost_price,omitempty"`
	Discount      float64 `gorm:"column:discount;type:decimal(5,2);default:100.00;comment:折扣(百分比)" json:"discount"`
	DiscountPrice float64 `gorm:"column:discount_price;type:decimal(10,2);default:null;comment:折后价" json:"discount_price,omitempty"`

	// 库存及销量相关字段
	Stock        int `gorm:"column:stock;type:int;default:0;comment:库存数量" json:"stock"`
	StockWarn    int `gorm:"column:stock_warn;type:int;default:10;comment:库存预警值" json:"stock_warn"`
	SoldCount    int `gorm:"column:sold_count;type:int;default:0;comment:已售数量" json:"sold_count"`
	ViewCount    int `gorm:"column:view_count;type:int;default:0;comment:浏览数量" json:"view_count"`
	CollectCount int `gorm:"column:collect_count;type:int;default:0;comment:收藏数量" json:"collect_count"`

	// 适用宠物相关字段
	PetType string  `gorm:"column:pet_type;type:enum('dog','cat','both','other');default:both;index;comment:宠物类型" json:"pet_type"`
	PetAge  []uint8 `gorm:"column:pet_age;type:json;default:null;comment:适用年龄[幼年,成年,老年]" json:"pet_age,omitempty"`
	PetSize []uint8 `gorm:"column:pet_size;type:json;default:null;comment:适用体型[小型,中型,大型]" json:"pet_size,omitempty"`

	// 商品规格相关字段
	Weight     float64 `gorm:"column:weight;type:decimal(8,2);default:null;comment:重量(kg)" json:"weight,omitempty"`
	Unit       string  `gorm:"column:unit;type:varchar(20);default:null;comment:单位(袋/盒/罐)" json:"unit,omitempty"`
	BrandID    uint64  `gorm:"column:brand_id;type:int;default:null;index;comment:品牌ID" json:"brand_id,omitempty"`
	CategoryID uint64  `gorm:"column:category_id;type:int;not null;index;comment:分类ID" json:"category_id"`
	Tags       []uint8 `gorm:"column:tags;type:json;default:null;comment:标签数组" json:"tags,omitempty"`

	// 商品详情相关字段
	NutritionInfo    []uint8 `gorm:"column:nutrition_info;type:json;default:null;comment:营养成分表" json:"nutrition_info,omitempty"` // 嵌套结构体对应复杂JSON
	Ingredients      string  `gorm:"column:ingredients;type:text;default:null;comment:主要成分" json:"ingredients,omitempty"`
	FeedingGuide     string  `gorm:"column:feeding_guide;type:text;default:null;comment:喂养指南" json:"feeding_guide,omitempty"`
	ShelfLife        int64   `gorm:"column:shelf_life;type:int;default:null;comment:保质期(天)" json:"shelf_life,omitempty"`
	StorageCondition string  `gorm:"column:storage_condition;type:varchar(200);default:null;comment:存储条件" json:"storage_condition,omitempty"`

	// 商品状态相关字段
	Status         string           `gorm:"column:status;type:enum('draft','pending','approved','rejected','on_sale','off_sale','deleted');default:draft;index;comment:状态" json:"status"`
	IsHot          bool             `gorm:"column:is_hot;type:tinyint(1);default:0;index;comment:是否热销" json:"is_hot"`
	IsNew          bool             `gorm:"column:is_new;type:tinyint(1);default:0;index;comment:是否新品" json:"is_new"`
	IsRecommend    bool             `gorm:"column:is_recommend;type:tinyint(1);default:0;index;comment:是否推荐" json:"is_recommend"`
	IsPrescription bool             `gorm:"column:is_prescription;type:tinyint(1);default:0;comment:是否处方粮" json:"is_prescription"`
	IsImported     bool             `gorm:"column:is_imported;type:tinyint(1);default:0;comment:是否进口" json:"is_imported"`
	IsOrganic      bool             `gorm:"column:is_organic;type:tinyint(1);default:0;comment:是否有机" json:"is_organic"`
	IsGrainFree    bool             `gorm:"column:is_grain_free;type:tinyint(1);default:0;comment:是否无谷" json:"is_grain_free"`
	PublishTime    common.LocalTime `gorm:"column:publish_time;type:datetime;default:null;comment:上架时间" json:"publish_time,omitempty"`

	DeletedAt common.LocalTime `gorm:"column:deleted_at;type:timestamp;default:null;index;comment:删除时间（软删除）" json:"deleted_at"`
}

type ProductListResp struct {
	// 基础核心字段（列表展示必备）
	ID        uint64 `gorm:"column:id" json:"id" comment:"商品ID"`
	ProductNo string `gorm:"column:product_no" json:"product_no" comment:"商品编号"`
	Name      string `gorm:"column:name" json:"name" comment:"商品名称"`
	Subtitle  string `gorm:"column:subtitle" json:"subtitle,omitempty" comment:"商品副标题"` // 可空字段，避免零值歧义
	MainImage string `gorm:"column:main_image" json:"main_image,omitempty" comment:"商品主图"`

	// 价格相关字段（列表展示必备，移除可计算的DiscountPrice）
	MarketPrice float64 `gorm:"column:market_price" json:"market_price,omitempty" comment:"市场价"`
	SalePrice   float64 `gorm:"column:sale_price" json:"sale_price" comment:"销售价"` // 非空必填，无需omitempty
	Discount    float64 `gorm:"column:discount" json:"discount" comment:"折扣(百分比，默认100.00)"`

	// 库存及销量相关字段（列表展示/库存预警所需）
	Stock     int `gorm:"column:stock" json:"stock" comment:"当前库存数量"`
	SoldCount int `gorm:"column:sold_count" json:"sold_count" comment:"累计已售数量"`

	// 分类与品牌字段（列表筛选/展示所需）
	CategoryID uint64 `gorm:"column:category_id" json:"category_id" comment:"分类ID"`
	BrandID    int64  `gorm:"column:brand_id" json:"brand_id,omitempty" comment:"品牌ID"`
	// 宠物适配字段（列表筛选所需）
	PetType string `gorm:"column:pet_type" json:"pet_type" comment:"适用宠物类型(dog/cat/both/other)"`

	// 商品标签字段（列表展示/筛选所需，JSON类型正确映射）
	Tags []uint8 `gorm:"column:tags;type:json" json:"tags,omitempty" comment:"商品标签数组（如\"爆款\",\"天然\"）"`

	// 商品状态字段（列表筛选/展示所需）
	Status         string `gorm:"column:status" json:"status" comment:"商品状态(on_sale/off_sale)"`
	IsHot          bool   `gorm:"column:is_hot" json:"is_hot" comment:"是否热销（首页推荐）"`
	IsNew          bool   `gorm:"column:is_new" json:"is_new" comment:"是否新品（新品专区）"`
	IsRecommend    bool   `gorm:"column:is_recommend" json:"is_recommend" comment:"是否推荐（精选列表）"`
	IsPrescription bool   `gorm:"column:is_prescription" json:"is_prescription,omitempty" comment:"是否处方粮（特殊标识）"`

	// 时间字段（列表排序/展示所需，指针类型标识可空）
	PublishTime common.LocalTime `gorm:"column:publish_time" json:"publish_time,omitempty" comment:"上架时间"`
	CreatedAt   common.LocalTime `gorm:"column:created_at" json:"created_at" comment:"创建时间"`
}

// NutritionItem 营养成分表嵌套结构体，对应nutrition_info字段的JSON结构
// 若营养成分表结构简单，也可直接使用map[string]float64，此处为强类型定义更规范
type NutritionItem struct {
	Name  string  `json:"name"`  // 营养成分名称（如：粗蛋白质）
	Value float64 `json:"value"` // 营养成分含量（如：25.5）
	Unit  string  `json:"unit"`  // 单位（如：%、g/100g）
}

// TableName 指定结构体对应的数据库表名
// GORM默认会将结构体名转为蛇形复数（如Product转为products），此处显式指定确保一致性
func (p *Product) TableName() string {
	return "products"
}
