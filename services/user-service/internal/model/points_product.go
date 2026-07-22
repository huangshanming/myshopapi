package model

import "mymall/common"

const (
	PointsProductStatusOn  = "on"
	PointsProductStatusOff = "off"
)

type PointsProduct struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	Name         string           `gorm:"column:name;type:varchar(100)" db:"name" json:"name"`
	CoverURL     string           `gorm:"column:cover_url;type:varchar(512)" db:"cover_url" json:"cover_url"`
	Description  string           `gorm:"column:description;type:varchar(1000)" db:"description" json:"description"`
	PointsPrice  int              `gorm:"column:points_price" db:"points_price" json:"points_price"`
	Stock        int              `gorm:"column:stock" db:"stock" json:"stock"`
	PerUserLimit int              `gorm:"column:per_user_limit" db:"per_user_limit" json:"per_user_limit"`
	Status       string           `gorm:"column:status;type:varchar(16)" db:"status" json:"status"`
	Sort         int              `gorm:"column:sort" db:"sort" json:"sort"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (PointsProduct) TableName() string { return "points_products" }
