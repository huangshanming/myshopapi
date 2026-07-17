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
	ID          uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID   uint64           `gorm:"column:product_id;not null;index" json:"product_id"`
	ShopID      uint64           `gorm:"column:shop_id;not null;index" json:"shop_id"`
	SkuNo       string           `gorm:"column:sku_no;type:varchar(64);uniqueIndex" json:"sku_no"`
	SpecValues  string           `gorm:"column:spec_values;type:json" json:"spec_values"` // JSON string
	SpecKey     string           `gorm:"column:spec_key;type:varchar(255)" json:"spec_key"`
	SalePrice   float64          `gorm:"column:sale_price;type:decimal(10,2)" json:"sale_price"`
	MarketPrice float64          `gorm:"column:market_price;type:decimal(10,2)" json:"market_price"`
	CostPrice   float64          `gorm:"column:cost_price;type:decimal(10,2)" json:"cost_price"`
	Stock       int              `gorm:"column:stock" json:"stock"`
	StockWarn   int              `gorm:"column:stock_warn" json:"stock_warn"`
	Barcode     string           `gorm:"column:barcode;type:varchar(64)" json:"barcode"`
	Status      string           `gorm:"column:status;type:enum('enabled','disabled');default:enabled" json:"status"`
	SoldCount   int              `gorm:"column:sold_count" json:"sold_count"`
	CreatedAt   common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *common.LocalTime `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

func (ProductSku) TableName() string { return "product_skus" }

type ProductImage struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID uint64           `gorm:"column:product_id;not null;index" json:"product_id"`
	ShopID    uint64           `gorm:"column:shop_id;not null" json:"shop_id"`
	URL       string           `gorm:"column:url;type:varchar(500)" json:"url"`
	Typ       string           `gorm:"column:typ;type:enum('main','gallery','detail');default:gallery" json:"typ"`
	Sort      int              `gorm:"column:sort" json:"sort"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ProductImage) TableName() string { return "product_images" }

type ProductTag struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID    uint64           `gorm:"column:shop_id;not null" json:"shop_id"`
	Name      string           `gorm:"column:name;type:varchar(50)" json:"name"`
	Color     string           `gorm:"column:color;type:varchar(20)" json:"color"`
	Status    int              `gorm:"column:status" json:"status"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ProductTag) TableName() string { return "product_tags" }

type ProductTagRel struct {
	ProductID uint64 `gorm:"column:product_id;primaryKey" json:"product_id"`
	TagID     uint64 `gorm:"column:tag_id;primaryKey" json:"tag_id"`
}

func (ProductTagRel) TableName() string { return "product_tag_rels" }

type ProductAttrTemplate struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID    uint64           `gorm:"column:shop_id;not null" json:"shop_id"`
	Name      string           `gorm:"column:name;type:varchar(100)" json:"name"`
	AttrsJSON string           `gorm:"column:attrs_json;type:json" json:"attrs_json"`
	Status    int              `gorm:"column:status" json:"status"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ProductAttrTemplate) TableName() string { return "product_attr_templates" }

type ProductAttr struct {
	ID         uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID  uint64 `gorm:"column:product_id;not null;index" json:"product_id"`
	TemplateID uint64 `gorm:"column:template_id" json:"template_id"`
	AttrKey    string `gorm:"column:attr_key;type:varchar(64)" json:"attr_key"`
	AttrLabel  string `gorm:"column:attr_label;type:varchar(100)" json:"attr_label"`
	AttrValue  string `gorm:"column:attr_value;type:varchar(500)" json:"attr_value"`
}

func (ProductAttr) TableName() string { return "product_attrs" }

type ProductSchedule struct {
	ID        uint64            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID uint64            `gorm:"column:product_id;not null" json:"product_id"`
	ShopID    uint64            `gorm:"column:shop_id;not null" json:"shop_id"`
	Action    string            `gorm:"column:action;type:enum('on_sale','off_sale')" json:"action"`
	RunAt     common.LocalTime  `gorm:"column:run_at" json:"run_at"`
	Status    string            `gorm:"column:status;type:enum('pending','done','cancelled');default:pending" json:"status"`
	LockedAt  *common.LocalTime `gorm:"column:locked_at" json:"locked_at,omitempty"`
	CreatedAt common.LocalTime  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime  `gorm:"column:updated_at" json:"updated_at"`
}

func (ProductSchedule) TableName() string { return "product_schedules" }

type ProductBatchJob struct {
	ID         uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID     uint64           `gorm:"column:shop_id;not null" json:"shop_id"`
	JobType    string           `gorm:"column:job_type;type:varchar(32)" json:"job_type"`
	PayloadJSON string          `gorm:"column:payload_json;type:json" json:"payload_json"`
	Progress   int              `gorm:"column:progress" json:"progress"`
	Total      int              `gorm:"column:total" json:"total"`
	Status     string           `gorm:"column:status" json:"status"`
	ResultMsg  string           `gorm:"column:result_msg;type:varchar(1000)" json:"result_msg"`
	OperatorID uint64           `gorm:"column:operator_id" json:"operator_id"`
	CreatedAt  common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ProductBatchJob) TableName() string { return "product_batch_jobs" }

type ProductOpLog struct {
	ID         uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID     uint64           `gorm:"column:shop_id;not null" json:"shop_id"`
	ProductID  *uint64          `gorm:"column:product_id" json:"product_id,omitempty"`
	OperatorID uint64           `gorm:"column:operator_id" json:"operator_id"`
	Action     string           `gorm:"column:action;type:varchar(64)" json:"action"`
	BeforeJSON string           `gorm:"column:before_json;type:json" json:"before_json,omitempty"`
	AfterJSON  string           `gorm:"column:after_json;type:json" json:"after_json,omitempty"`
	IP         string           `gorm:"column:ip;type:varchar(64)" json:"ip"`
	CreatedAt  common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (ProductOpLog) TableName() string { return "product_op_logs" }

// 商家 RBAC
type ShopRole struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID    uint64           `gorm:"column:shop_id;not null" json:"shop_id"`
	Code      string           `gorm:"column:code;type:varchar(64)" json:"code"`
	Name      string           `gorm:"column:name;type:varchar(100)" json:"name"`
	Status    int              `gorm:"column:status" json:"status"`
	Remark    string           `gorm:"column:remark;type:varchar(255)" json:"remark"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ShopRole) TableName() string { return "shop_roles" }

type ShopMenu struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentID  uint64           `gorm:"column:parent_id" json:"parent_id"`
	Name      string           `gorm:"column:name" json:"name"`
	Type      string           `gorm:"column:type" json:"type"`
	Path      string           `gorm:"column:path" json:"path"`
	Component string           `gorm:"column:component" json:"component"`
	Icon      string           `gorm:"column:icon" json:"icon"`
	Perms     string           `gorm:"column:perms" json:"perms"`
	Sort      int              `gorm:"column:sort" json:"sort"`
	Visible   int              `gorm:"column:visible" json:"visible"`
	Status    int              `gorm:"column:status" json:"status"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ShopMenu) TableName() string { return "shop_menus" }

type ShopRoleMenu struct {
	RoleID uint64 `gorm:"column:role_id;primaryKey" json:"role_id"`
	MenuID uint64 `gorm:"column:menu_id;primaryKey" json:"menu_id"`
}

func (ShopRoleMenu) TableName() string { return "shop_role_menus" }

type ShopUserRole struct {
	ShopID uint64 `gorm:"column:shop_id;primaryKey" json:"shop_id"`
	UserID uint64 `gorm:"column:user_id;primaryKey" json:"user_id"`
	RoleID uint64 `gorm:"column:role_id;primaryKey" json:"role_id"`
}

func (ShopUserRole) TableName() string { return "shop_user_roles" }
