package model

import "mymall/common"

const (
	ReviewStatusVisible = "visible"
	ReviewStatusDeleted = "deleted"
)

type ProductReview struct {
	ID            uint64               `db:"id" json:"id"`
	OrderID       uint64               `db:"order_id" json:"order_id"`
	OrderNo       string               `db:"order_no" json:"order_no"`
	UserID        uint64               `db:"user_id" json:"user_id"`
	ShopID        uint64               `db:"shop_id" json:"shop_id"`
	ProductID     uint64               `db:"product_id" json:"product_id"`
	OrderItemID   uint64               `db:"order_item_id" json:"order_item_id"`
	SkuID         uint64               `db:"sku_id" json:"sku_id"`
	SkuSnapshot   string               `db:"sku_snapshot" json:"sku_snapshot"`
	Rating        int8                 `db:"rating" json:"rating"`
	Content       string               `db:"content" json:"content"`
	IsAnonymous   bool                 `db:"is_anonymous" json:"is_anonymous"`
	Status        string               `db:"status" json:"status"`
	MerchantReply string               `db:"merchant_reply" json:"merchant_reply"`
	RepliedAt     *common.LocalTime    `db:"replied_at" json:"replied_at,omitempty"`
	CreatedAt     common.LocalTime     `db:"created_at" json:"created_at"`
	UpdatedAt     common.LocalTime     `db:"updated_at" json:"updated_at"`
	Images        []ProductReviewImage `db:"-" json:"images,omitempty"`
	UserName      string               `db:"-" json:"user_name,omitempty"`
}

func (ProductReview) TableName() string { return "product_reviews" }

type ProductReviewImage struct {
	ID       uint64 `db:"id" json:"id"`
	ReviewID uint64 `db:"review_id" json:"review_id"`
	URL      string `db:"url" json:"url"`
	Sort     int    `db:"sort" json:"sort"`
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
