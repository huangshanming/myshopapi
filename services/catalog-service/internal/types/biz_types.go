package types

// Hand-maintained API DTOs for catalog-service (sync with api/*.api).
// Nested product/article shapes mirror domain packages; logic converts when needed.

type SpecItem struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type SkuInput struct {
	ID          uint64            `json:"id,optional"`
	SpecValues  map[string]string `json:"spec_values,optional"`
	SalePrice   float64           `json:"sale_price,optional"`
	MarketPrice float64           `json:"market_price,optional"`
	CostPrice   float64           `json:"cost_price,optional"`
	Stock       int               `json:"stock,optional"`
	StockWarn   int               `json:"stock_warn,optional"`
	Barcode     string            `json:"barcode,optional"`
	Status      string            `json:"status,optional"`
}

type ImageInput struct {
	URL  string `json:"url"`
	Typ  string `json:"typ,optional"`
	Sort int    `json:"sort,optional"`
}

type AttrInput struct {
	TemplateID uint64 `json:"template_id,optional"`
	AttrKey    string `json:"attr_key,optional"`
	AttrLabel  string `json:"attr_label,optional"`
	AttrValue  string `json:"attr_value,optional"`
}

type MerchantProductSaveReq struct {
	Name             string       `json:"name"`
	Subtitle         string       `json:"subtitle,optional"`
	Description      string       `json:"description,optional"`
	CategoryID       uint64       `json:"category_id,optional"`
	ProductType      string       `json:"product_type,optional"`
	PetType          string       `json:"pet_type,optional"`
	Status           string       `json:"status,optional"`
	MainImage        string       `json:"main_image,optional"`
	SpecJSON         []SpecItem   `json:"spec_json,optional"`
	Skus             []SkuInput   `json:"skus,optional"`
	Images           []ImageInput `json:"images,optional"`
	Attrs            []AttrInput  `json:"attrs,optional"`
	TagIDs           []uint64     `json:"tag_ids,optional"`
	ShelfLife        int64        `json:"shelf_life,optional"`
	StorageCondition string       `json:"storage_condition,optional"`
	SalePrice        float64      `json:"sale_price,optional"`
	Stock            int          `json:"stock,optional"`
}

type BatchProductReq struct {
	Action     string   `json:"action"`
	ProductIDs []uint64 `json:"product_ids"`
	CategoryID uint64   `json:"category_id,optional"`
	PriceDelta float64  `json:"price_delta,optional"`
	PriceRate  float64  `json:"price_rate,optional"`
}

type StockAdjustReq struct {
	SkuID uint64 `json:"sku_id"`
	Stock *int   `json:"stock,optional"`
	Delta *int   `json:"delta,optional"`
}

type BatchStockReq struct {
	Items []StockAdjustReq `json:"items"`
}

type RecycleReq struct {
	ProductIDs []uint64 `json:"product_ids"`
}

type TagReq struct {
	Name  string `json:"name"`
	Color string `json:"color,optional"`
}

type AttrTemplateReq struct {
	Name      string `json:"name"`
	AttrsJSON string `json:"attrs_json,optional"`
}

type CategoryReq struct {
	ParentId    uint64 `json:"parent_id,optional"`
	Name        string `json:"name"`
	Icon        string `json:"icon,optional"`
	Description string `json:"description,optional"`
	SortOrder   int    `json:"sort_order,optional"`
	Level       int    `json:"level,optional"`
	IsShow      *bool  `json:"is_show,optional"`
}

type PlatformProductRemarkReq struct {
	Remark string `json:"remark,optional"`
}

type ArticleSaveReq struct {
	ShopID            uint64   `json:"shop_id,optional"`
	CategoryID        uint64   `json:"category_id,optional"`
	Title             string   `json:"title"`
	CoverURL          string   `json:"cover_url,optional"`
	Content           string   `json:"content,optional"`
	SchedulePublishAt string   `json:"schedule_publish_at,optional"`
	IsTop             *int8    `json:"is_top,optional"`
	ImageURLs         []string `json:"image_urls,optional"`
	Submit            bool     `json:"submit,optional"`
}

type ArticleBatchAuditReq struct {
	IDs          []uint64 `json:"ids"`
	Pass         bool     `json:"pass"`
	RejectReason string   `json:"reject_reason,optional"`
}

type ArticleCategorySaveReq struct {
	ParentID uint64 `json:"parent_id,optional"`
	Name     string `json:"name"`
	Sort     int    `json:"sort,optional"`
	Status   *int8  `json:"status,optional"`
}

type ArticleIdListReq struct {
	Id uint64 `json:"id"`
}

type FavoriteAddReq struct {
	ProductID uint64 `json:"product_id"`
}

type FavoriteBatchRemoveReq struct {
	ProductIDs []uint64 `json:"product_ids"`
}

type UserArticleCreateReq struct {
	CategoryID uint64   `json:"category_id,optional"`
	Title      string   `json:"title"`
	CoverURL   string   `json:"cover_url,optional"`
	Content    string   `json:"content,optional"`
	ImageURLs  []string `json:"image_urls,optional"`
}

type ShopRoleReq struct {
	Code    string   `json:"code"`
	Name    string   `json:"name"`
	Remark  string   `json:"remark,optional"`
	MenuIDs []uint64 `json:"menu_ids,optional"`
}

type ShopStaffReq struct {
	Mobile   string `json:"mobile"`
	RoleID   uint64 `json:"role_id"`
	Nickname string `json:"nickname,optional"`
	Password string `json:"password,optional"`
	Mode     string `json:"mode,optional"`
}

type BannerSaveReq struct {
	Title    string `json:"title,optional"`
	ImageURL string `json:"image_url"`
	LinkType string `json:"link_type,optional"`
	LinkID   uint64 `json:"link_id,optional"`
	Sort     int    `json:"sort,optional"`
	Status   string `json:"status,optional"`
	StartAt  string `json:"start_at,optional"`
	EndAt    string `json:"end_at,optional"`
}

type EmojiSaveReq struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Sort     int    `json:"sort,optional"`
	Status   *int8  `json:"status,optional"`
}

// ---- path + body ----

type ProductUpdateBodyReq struct {
	Id uint64 `path:"id"`
	MerchantProductSaveReq
}

type SetStatusBodyReq struct {
	Id     uint64 `path:"id"`
	Status string `json:"status"`
}

type TagUpdateBodyReq struct {
	Id    uint64 `path:"id"`
	Name  string `json:"name"`
	Color string `json:"color,optional"`
}

type AttrTemplateUpdateBodyReq struct {
	Id        uint64 `path:"id"`
	Name      string `json:"name"`
	AttrsJSON string `json:"attrs_json,optional"`
}

type StockAdjustBodyReq struct {
	Id    uint64 `path:"id"`
	SkuID uint64 `json:"sku_id"`
	Stock *int   `json:"stock,optional"`
	Delta *int   `json:"delta,optional"`
}

type ScheduleBodyReq struct {
	Id     uint64 `path:"id"`
	Action string `json:"action"`
	RunAt  string `json:"run_at"`
}

type ShopRoleUpdateBodyReq struct {
	Id      uint64   `path:"id"`
	Code    string   `json:"code"`
	Name    string   `json:"name"`
	Remark  string   `json:"remark,optional"`
	MenuIDs []uint64 `json:"menu_ids,optional"`
}

type ArticleUpdateBodyReq struct {
	Id uint64 `path:"id"`
	ArticleSaveReq
}

type ArticleCommentPatchBodyReq struct {
	Id     uint64 `path:"id"`
	Status string `json:"status"`
}

type ArticleCategoryUpdateBodyReq struct {
	Id       uint64 `path:"id"`
	ParentID uint64 `json:"parent_id,optional"`
	Name     string `json:"name"`
	Sort     int    `json:"sort,optional"`
	Status   *int8  `json:"status,optional"`
}

type ArticleTopBodyReq struct {
	Id    uint64 `path:"id"`
	IsTop int8   `json:"is_top"`
}

type ArticleRemarkBodyReq struct {
	Id     uint64 `path:"id"`
	Remark string `json:"remark,optional"`
}

type ArticleAuditBodyReq struct {
	Id           uint64 `path:"id"`
	Pass         bool   `json:"pass"`
	RejectReason string `json:"reject_reason,optional"`
}

type UserArticleUpdateBodyReq struct {
	Id         uint64   `path:"id"`
	CategoryID uint64   `json:"category_id,optional"`
	Title      string   `json:"title"`
	CoverURL   string   `json:"cover_url,optional"`
	Content    string   `json:"content,optional"`
	ImageURLs  []string `json:"image_urls,optional"`
}

type CreateCommentBodyReq struct {
	Id       uint64 `path:"id"`
	Content  string `json:"content"`
	ParentID uint64 `json:"parent_id,optional"`
}

type BannerUpdateBodyReq struct {
	Id       uint64 `path:"id"`
	Title    string `json:"title,optional"`
	ImageURL string `json:"image_url"`
	LinkType string `json:"link_type,optional"`
	LinkID   uint64 `json:"link_id,optional"`
	Sort     int    `json:"sort,optional"`
	Status   string `json:"status,optional"`
	StartAt  string `json:"start_at,optional"`
	EndAt    string `json:"end_at,optional"`
}

type PlatformProductRemarkBodyReq struct {
	Id     uint64 `path:"id"`
	Remark string `json:"remark,optional"`
}

type CategoryUpdateBodyReq struct {
	Id          uint64 `path:"id"`
	ParentId    uint64 `json:"parent_id,optional"`
	Name        string `json:"name"`
	Icon        string `json:"icon,optional"`
	Description string `json:"description,optional"`
	SortOrder   int    `json:"sort_order,optional"`
	Level       int    `json:"level,optional"`
	IsShow      *bool  `json:"is_show,optional"`
}

type EmojiUpdateBodyReq struct {
	Id       uint64 `path:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Sort     int    `json:"sort,optional"`
	Status   *int8  `json:"status,optional"`
}

type SlotTypeQueryReq struct {
	SlotType string `form:"slot_type,optional"`
}

type IdQueryReq struct {
	Id uint64 `form:"id"`
}

type ProductIdPathReq struct {
	ProductId uint64 `path:"product_id"`
}

// ---- list / query (go-zero form tags) ----

type IdPageReq struct {
	Id       uint64 `path:"id"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type PublicProductListReq struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	ShopId     uint64 `form:"shop_id,optional"`
	CategoryId uint64 `form:"category_id,optional"`
	OrderBy    string `form:"order_by,optional"`
}

type PublicArticleListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Home     string `form:"home,optional"`
}

type MerchantProductListReq struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	Name        string `form:"name,optional"`
	ProductNo   string `form:"product_no,optional"`
	CategoryId  uint64 `form:"category_id,optional"`
	Status      string `form:"status,optional"`
	ProductType string `form:"product_type,optional"`
	StockWarn   string `form:"stock_warn,optional"`
	OrderBy     string `form:"order_by,optional"`
	Recycle     string `form:"recycle,optional"`
}

type AdminProductListReq struct {
	Page         int    `form:"page,default=1"`
	PageSize     int    `form:"page_size,default=10"`
	ShopId       uint64 `form:"shop_id,optional"`
	Name         string `form:"name,optional"`
	ProductNo    string `form:"product_no,optional"`
	CategoryId   uint64 `form:"category_id,optional"`
	Status       string `form:"status,optional"`
	ProductType  string `form:"product_type,optional"`
	OrderBy      string `form:"order_by,optional"`
	CreatedFrom  string `form:"created_from,optional"`
	CreatedTo    string `form:"created_to,optional"`
	PublishFrom  string `form:"publish_from,optional"`
	PublishTo    string `form:"publish_to,optional"`
}

type MerchantArticleListReq struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	Title       string `form:"title,optional"`
	AuditStatus string `form:"audit_status,optional"`
	Status      string `form:"status,optional"`
}

type AdminArticleListReq struct {
	Page         int    `form:"page,default=1"`
	PageSize     int    `form:"page_size,default=10"`
	Title        string `form:"title,optional"`
	AuditStatus  string `form:"audit_status,optional"`
	Status       string `form:"status,optional"`
	ShopId       uint64 `form:"shop_id,optional"`
	HasSchedule  string `form:"has_schedule,optional"`
	CreatedFrom  string `form:"created_from,optional"`
	CreatedTo    string `form:"created_to,optional"`
}

type ArticleCommentListReq struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	ArticleId uint64 `form:"article_id,optional"`
	ShopId    uint64 `form:"shop_id,optional"`
	Status    string `form:"status,optional"`
}

type NotificationListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	IsRead   string `form:"is_read,optional"`
}

type OpLogsReq struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	ProductId uint64 `form:"product_id,optional"`
}
