package model

import "mymall/common"

type ProductFavorite struct {
	ID        uint64           `db:"id" json:"id"`
	UserID    uint64           `db:"user_id" json:"user_id"`
	ProductID uint64           `db:"product_id" json:"product_id"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
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
