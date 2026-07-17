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
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID    uint64           `gorm:"column:shop_id;not null;index" json:"shop_id"`
	Type      string           `gorm:"column:type;type:varchar(32);not null" json:"type"`
	Title     string           `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Content   string           `gorm:"column:content;type:varchar(1000)" json:"content"`
	Link      string           `gorm:"column:link;type:varchar(255)" json:"link"`
	RefType   string           `gorm:"column:ref_type;type:varchar(32)" json:"ref_type"`
	RefID     uint64           `gorm:"column:ref_id;default:0" json:"ref_id"`
	IsRead    int8             `gorm:"column:is_read;default:0" json:"is_read"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (ShopNotification) TableName() string { return "shop_notifications" }
