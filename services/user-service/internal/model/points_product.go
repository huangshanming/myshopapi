package model

import "mymall/common"

const (
	PointsProductStatusOn  = "on"
	PointsProductStatusOff = "off"
)

type PointsProduct struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string           `gorm:"column:name;type:varchar(100)" json:"name"`
	CoverURL     string           `gorm:"column:cover_url;type:varchar(512)" json:"cover_url"`
	Description  string           `gorm:"column:description;type:varchar(1000)" json:"description"`
	PointsPrice  int              `gorm:"column:points_price" json:"points_price"`
	Stock        int              `gorm:"column:stock" json:"stock"`
	PerUserLimit int              `gorm:"column:per_user_limit" json:"per_user_limit"`
	Status       string           `gorm:"column:status;type:varchar(16)" json:"status"`
	Sort         int              `gorm:"column:sort" json:"sort"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (PointsProduct) TableName() string { return "points_products" }
