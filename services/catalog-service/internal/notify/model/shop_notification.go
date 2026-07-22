package model

import "mymall/common"

const (
	NotifProductOffSale  = "product_off_sale"
	NotifProductDeleted  = "product_deleted"
	NotifArticleOffline  = "article_offline"
	NotifArticleDeleted  = "article_deleted"
	NotifArticleRejected = "article_rejected"
)

type ShopNotification struct {
	ID        uint64           `db:"id" json:"id"`
	ShopID    uint64           `db:"shop_id" json:"shop_id"`
	Type      string           `db:"type" json:"type"`
	Title     string           `db:"title" json:"title"`
	Content   string           `db:"content" json:"content"`
	Link      string           `db:"link" json:"link"`
	RefType   string           `db:"ref_type" json:"ref_type"`
	RefID     uint64           `db:"ref_id" json:"ref_id"`
	IsRead    int8             `db:"is_read" json:"is_read"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
}

func (ShopNotification) TableName() string { return "shop_notifications" }
