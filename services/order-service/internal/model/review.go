package model

import "mymall/common"

const (
	ReviewStatusVisible = "visible"
	ReviewStatusDeleted = "deleted"
)

type ProductReview struct {
	ID            uint64            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID       uint64            `gorm:"column:order_id;not null;uniqueIndex" json:"order_id"`
	OrderNo       string            `gorm:"column:order_no;type:varchar(64);not null" json:"order_no"`
	UserID        uint64            `gorm:"column:user_id;not null;index" json:"user_id"`
	ShopID        uint64            `gorm:"column:shop_id;not null;default:0;index" json:"shop_id"`
	ProductID     uint64            `gorm:"column:product_id;not null;index" json:"product_id"`
	OrderItemID   uint64            `gorm:"column:order_item_id;not null;default:0" json:"order_item_id"`
	SkuID         uint64            `gorm:"column:sku_id;not null;default:0" json:"sku_id"`
	SkuSnapshot   string            `gorm:"column:sku_snapshot;type:varchar(512);not null;default:''" json:"sku_snapshot"`
	Rating        int8              `gorm:"column:rating;not null" json:"rating"`
	Content       string            `gorm:"column:content;type:varchar(1000);not null;default:''" json:"content"`
	IsAnonymous   bool              `gorm:"column:is_anonymous;type:tinyint;not null;default:0" json:"is_anonymous"`
	Status        string            `gorm:"column:status;type:varchar(16);not null;default:visible" json:"status"`
	MerchantReply string            `gorm:"column:merchant_reply;type:varchar(500);not null;default:''" json:"merchant_reply"`
	RepliedAt     *common.LocalTime `gorm:"column:replied_at" json:"replied_at,omitempty"`
	CreatedAt     common.LocalTime  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     common.LocalTime  `gorm:"column:updated_at" json:"updated_at"`

	Images   []ProductReviewImage `gorm:"foreignKey:ReviewID" json:"images,omitempty"`
	UserName string               `gorm:"-" json:"user_name,omitempty"`
}

func (ProductReview) TableName() string { return "product_reviews" }

type ProductReviewImage struct {
	ID       uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ReviewID uint64 `gorm:"column:review_id;not null;index" json:"review_id"`
	URL      string `gorm:"column:url;type:varchar(500);not null" json:"url"`
	Sort     int    `gorm:"column:sort;not null;default:0" json:"sort"`
}

func (ProductReviewImage) TableName() string { return "product_review_images" }

type CreateReviewReq struct {
	Rating      int8     `json:"rating"`
	Content     string   `json:"content"`
	IsAnonymous bool     `json:"is_anonymous"`
	OrderItemID uint64   `json:"order_item_id"`
	Images      []string `json:"images"`
}

type ReviewListItem struct {
	ProductReview
	ProductName string `json:"product_name,omitempty"`
}
