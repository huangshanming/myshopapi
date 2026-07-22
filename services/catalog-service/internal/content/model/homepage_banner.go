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
	ID        uint64            `db:"id" json:"id"`
	Title     string            `db:"title" json:"title"`
	ImageURL  string            `db:"image_url" json:"image_url"`
	LinkType  string            `db:"link_type" json:"link_type"`
	LinkID    uint64            `db:"link_id" json:"link_id"`
	Sort      int               `db:"sort" json:"sort"`
	Status    string            `db:"status" json:"status"`
	StartAt   *common.LocalTime `db:"start_at" json:"start_at,omitempty"`
	EndAt     *common.LocalTime `db:"end_at" json:"end_at,omitempty"`
	CreatedAt common.LocalTime  `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime  `db:"updated_at" json:"updated_at"`
	// 展示字段
	LinkName string `db:"-" json:"link_name,omitempty"`
}

func (HomepageBanner) TableName() string { return "homepage_banners" }
