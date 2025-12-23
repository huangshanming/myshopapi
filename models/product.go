package model

import (
	"database/sql"
	"time"
)

// Product 商品主表模型，对应数据库表products
// 字段类型严格匹配数据库表定义，标签包含GORM配置、JSON序列化配置及注释
type Product struct {
	// 基础字段
	ID          uint64         `gorm:"column:id;primaryKey;autoIncrement;comment:商品ID" json:"id"`
	ProductNo   string         `gorm:"column:product_no;type:varchar(50);not null;uniqueIndex;comment:商品编号" json:"product_no"`
	Name        string         `gorm:"column:name;type:varchar(200);not null;comment:商品名称" json:"name"`
	Subtitle    sql.NullString `gorm:"column:subtitle;type:varchar(500);default:null;comment:商品副标题" json:"subtitle,omitempty"` // 可空字段使用sql.NullString避免空指针
	Description sql.NullString `gorm:"column:description;type:text;default:null;comment:商品描述" json:"description,omitempty"`
	MainImage   sql.NullString `gorm:"column:main_image;type:varchar(500);default:null;comment:主图" json:"main_image,omitempty"`
	ImageList   []uint8        `gorm:"column:image_list;type:json;default:null;comment:商品图片列表" json:"image_list,omitempty"` // JSON数组映射为string切片
	VideoUrl    sql.NullString `gorm:"column:video_url;type:varchar(500);default:null;comment:商品视频" json:"video_url,omitempty"`

	// 价格相关字段
	MarketPrice   sql.NullFloat64 `gorm:"column:market_price;type:decimal(10,2);default:null;comment:市场价" json:"market_price,omitempty"`
	SalePrice     float64         `gorm:"column:sale_price;type:decimal(10,2);not null;comment:销售价" json:"sale_price"`
	CostPrice     sql.NullFloat64 `gorm:"column:cost_price;type:decimal(10,2);default:null;comment:成本价" json:"cost_price,omitempty"`
	Discount      float64         `gorm:"column:discount;type:decimal(5,2);default:100.00;comment:折扣(百分比)" json:"discount"`
	DiscountPrice sql.NullFloat64 `gorm:"column:discount_price;type:decimal(10,2);default:null;comment:折后价" json:"discount_price,omitempty"`

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
	Weight     sql.NullFloat64 `gorm:"column:weight;type:decimal(8,2);default:null;comment:重量(kg)" json:"weight,omitempty"`
	Unit       sql.NullString  `gorm:"column:unit;type:varchar(20);default:null;comment:单位(袋/盒/罐)" json:"unit,omitempty"`
	BrandID    sql.NullInt64   `gorm:"column:brand_id;type:int;default:null;index;comment:品牌ID" json:"brand_id,omitempty"`
	CategoryID uint64          `gorm:"column:category_id;type:int;not null;index;comment:分类ID" json:"category_id"`
	Tags       []uint8         `gorm:"column:tags;type:json;default:null;comment:标签数组" json:"tags,omitempty"`

	// 商品详情相关字段
	NutritionInfo    []uint8        `gorm:"column:nutrition_info;type:json;default:null;comment:营养成分表" json:"nutrition_info,omitempty"` // 嵌套结构体对应复杂JSON
	Ingredients      sql.NullString `gorm:"column:ingredients;type:text;default:null;comment:主要成分" json:"ingredients,omitempty"`
	FeedingGuide     sql.NullString `gorm:"column:feeding_guide;type:text;default:null;comment:喂养指南" json:"feeding_guide,omitempty"`
	ShelfLife        sql.NullInt64  `gorm:"column:shelf_life;type:int;default:null;comment:保质期(天)" json:"shelf_life,omitempty"`
	StorageCondition sql.NullString `gorm:"column:storage_condition;type:varchar(200);default:null;comment:存储条件" json:"storage_condition,omitempty"`

	// 商品状态相关字段
	Status         string       `gorm:"column:status;type:enum('draft','pending','approved','rejected','on_sale','off_sale','deleted');default:draft;index;comment:状态" json:"status"`
	IsHot          bool         `gorm:"column:is_hot;type:tinyint(1);default:0;index;comment:是否热销" json:"is_hot"`
	IsNew          bool         `gorm:"column:is_new;type:tinyint(1);default:0;index;comment:是否新品" json:"is_new"`
	IsRecommend    bool         `gorm:"column:is_recommend;type:tinyint(1);default:0;index;comment:是否推荐" json:"is_recommend"`
	IsPrescription bool         `gorm:"column:is_prescription;type:tinyint(1);default:0;comment:是否处方粮" json:"is_prescription"`
	IsImported     bool         `gorm:"column:is_imported;type:tinyint(1);default:0;comment:是否进口" json:"is_imported"`
	IsOrganic      bool         `gorm:"column:is_organic;type:tinyint(1);default:0;comment:是否有机" json:"is_organic"`
	IsGrainFree    bool         `gorm:"column:is_grain_free;type:tinyint(1);default:0;comment:是否无谷" json:"is_grain_free"`
	PublishTime    sql.NullTime `gorm:"column:publish_time;type:datetime;default:null;comment:上架时间" json:"publish_time,omitempty"`

	// 时间戳字段（软删除字段deleted_at需配合GORM软删除功能）
	CreatedAt time.Time    `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time    `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at;type:timestamp;default:null;index;comment:删除时间（软删除）" json:"deleted_at,omitempty"`
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
