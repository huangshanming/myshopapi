package model

import "mymall/common"

// 商品状态（商家中台）
const (
	ProductDraft   = "draft"
	ProductOnSale  = "on_sale"
	ProductOffSale = "off_sale"
	ProductDeleted = "deleted"
)

const (
	ProductTypePhysical = "physical"
	ProductTypeFresh    = "fresh"
	ProductTypeVirtual  = "virtual"
)

const (
	SKUEnabled  = "enabled"
	SKUDisabled = "disabled"
)

// ProductSku SKU 规格库存
type ProductSku struct {
	ID          uint64            `db:"id" json:"id"`
	ProductID   uint64            `db:"product_id" json:"product_id"`
	ShopID      uint64            `db:"shop_id" json:"shop_id"`
	SkuNo       string            `db:"sku_no" json:"sku_no"`
	SpecValues  string            `db:"spec_values" json:"spec_values"` // JSON string
	SpecKey     string            `db:"spec_key" json:"spec_key"`
	SalePrice   float64           `db:"sale_price" json:"sale_price"`
	MarketPrice float64           `db:"market_price" json:"market_price"`
	CostPrice   float64           `db:"cost_price" json:"cost_price"`
	Stock       int               `db:"stock" json:"stock"`
	StockWarn   int               `db:"stock_warn" json:"stock_warn"`
	Barcode     string            `db:"barcode" json:"barcode"`
	Status      string            `db:"status" json:"status"`
	SoldCount   int               `db:"sold_count" json:"sold_count"`
	CreatedAt   common.LocalTime  `db:"created_at" json:"created_at"`
	UpdatedAt   common.LocalTime  `db:"updated_at" json:"updated_at"`
	DeletedAt   *common.LocalTime `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (ProductSku) TableName() string { return "product_skus" }

type ProductImage struct {
	ID        uint64           `db:"id" json:"id"`
	ProductID uint64           `db:"product_id" json:"product_id"`
	ShopID    uint64           `db:"shop_id" json:"shop_id"`
	URL       string           `db:"url" json:"url"`
	Typ       string           `db:"typ" json:"typ"`
	Sort      int              `db:"sort" json:"sort"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (ProductImage) TableName() string { return "product_images" }

type ProductTag struct {
	ID        uint64           `db:"id" json:"id"`
	ShopID    uint64           `db:"shop_id" json:"shop_id"`
	Name      string           `db:"name" json:"name"`
	Color     string           `db:"color" json:"color"`
	Status    int              `db:"status" json:"status"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (ProductTag) TableName() string { return "product_tags" }

type ProductTagRel struct {
	ProductID uint64 `db:"product_id" json:"product_id"`
	TagID     uint64 `db:"tag_id" json:"tag_id"`
}

func (ProductTagRel) TableName() string { return "product_tag_rels" }

type ProductAttrTemplate struct {
	ID        uint64           `db:"id" json:"id"`
	ShopID    uint64           `db:"shop_id" json:"shop_id"`
	Name      string           `db:"name" json:"name"`
	AttrsJSON string           `db:"attrs_json" json:"attrs_json"`
	Status    int              `db:"status" json:"status"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (ProductAttrTemplate) TableName() string { return "product_attr_templates" }

type ProductAttr struct {
	ID         uint64 `db:"id" json:"id"`
	ProductID  uint64 `db:"product_id" json:"product_id"`
	TemplateID uint64 `db:"template_id" json:"template_id"`
	AttrKey    string `db:"attr_key" json:"attr_key"`
	AttrLabel  string `db:"attr_label" json:"attr_label"`
	AttrValue  string `db:"attr_value" json:"attr_value"`
}

func (ProductAttr) TableName() string { return "product_attrs" }

type ProductSchedule struct {
	ID        uint64            `db:"id" json:"id"`
	ProductID uint64            `db:"product_id" json:"product_id"`
	ShopID    uint64            `db:"shop_id" json:"shop_id"`
	Action    string            `db:"action" json:"action"`
	RunAt     common.LocalTime  `db:"run_at" json:"run_at"`
	Status    string            `db:"status" json:"status"`
	LockedAt  *common.LocalTime `db:"locked_at" json:"locked_at,omitempty"`
	CreatedAt common.LocalTime  `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime  `db:"updated_at" json:"updated_at"`
}

func (ProductSchedule) TableName() string { return "product_schedules" }

type ProductBatchJob struct {
	ID          uint64           `db:"id" json:"id"`
	ShopID      uint64           `db:"shop_id" json:"shop_id"`
	JobType     string           `db:"job_type" json:"job_type"`
	PayloadJSON string           `db:"payload_json" json:"payload_json"`
	Progress    int              `db:"progress" json:"progress"`
	Total       int              `db:"total" json:"total"`
	Status      string           `db:"status" json:"status"`
	ResultMsg   string           `db:"result_msg" json:"result_msg"`
	OperatorID  uint64           `db:"operator_id" json:"operator_id"`
	CreatedAt   common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt   common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (ProductBatchJob) TableName() string { return "product_batch_jobs" }

type ProductOpLog struct {
	ID         uint64           `db:"id" json:"id"`
	ShopID     uint64           `db:"shop_id" json:"shop_id"`
	// Non-pointer: go-zero sqlx pre-allocates *uint64 so NULL cannot scan; use IFNULL in SELECT.
	ProductID  uint64           `db:"product_id" json:"product_id,omitempty"`
	OperatorID uint64           `db:"operator_id" json:"operator_id"`
	Action     string           `db:"action" json:"action"`
	BeforeJSON string           `db:"before_json" json:"before_json,omitempty"`
	AfterJSON  string           `db:"after_json" json:"after_json,omitempty"`
	IP         string           `db:"ip" json:"ip"`
	CreatedAt  common.LocalTime `db:"created_at" json:"created_at"`
}

func (ProductOpLog) TableName() string { return "product_op_logs" }
