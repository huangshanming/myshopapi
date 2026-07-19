package model

import "mymall/common"

type ProductFavorite struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint64           `gorm:"column:user_id;not null;uniqueIndex:uk_user_product" json:"user_id"`
	ProductID uint64           `gorm:"column:product_id;not null;uniqueIndex:uk_user_product;index" json:"product_id"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (ProductFavorite) TableName() string { return "product_favorites" }

type FavoriteListItem struct {
	ID           uint64           `json:"id"`
	ProductID    uint64           `json:"product_id"`
	Name         string           `json:"name"`
	MainImage    string           `json:"main_image"`
	SalePrice    float64          `json:"sale_price"`
	Status       string           `json:"status"`
	Invalid      bool             `json:"invalid"`
	CollectCount int              `json:"collect_count"`
	CreatedAt    common.LocalTime `json:"created_at"`
}
