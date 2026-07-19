package model

import "mymall/common"

const (
	BannerLinkNone    = "none"
	BannerLinkProduct = "product"
	BannerLinkArticle = "article"
)

const (
	BannerOn  = "on"
	BannerOff = "off"
)

type HomepageBanner struct {
	ID        uint64            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title     string            `gorm:"column:title;type:varchar(100);not null;default:''" json:"title"`
	ImageURL  string            `gorm:"column:image_url;type:varchar(500);not null" json:"image_url"`
	LinkType  string            `gorm:"column:link_type;type:enum('none','product','article');default:none" json:"link_type"`
	LinkID    uint64            `gorm:"column:link_id;not null;default:0" json:"link_id"`
	Sort      int               `gorm:"column:sort;default:0" json:"sort"`
	Status    string            `gorm:"column:status;type:enum('on','off');default:on" json:"status"`
	StartAt   *common.LocalTime `gorm:"column:start_at" json:"start_at,omitempty"`
	EndAt     *common.LocalTime `gorm:"column:end_at" json:"end_at,omitempty"`
	CreatedAt common.LocalTime  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime  `gorm:"column:updated_at" json:"updated_at"`
	// 展示字段
	LinkName string `gorm:"-" json:"link_name,omitempty"`
}

func (HomepageBanner) TableName() string { return "homepage_banners" }
