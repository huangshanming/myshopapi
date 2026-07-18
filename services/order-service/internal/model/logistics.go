package model

import "mymall/common"

type LogisticsCompany struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string           `gorm:"column:name;type:varchar(64);not null" json:"name"`
	Code      string           `gorm:"column:code;type:varchar(32);not null;uniqueIndex" json:"code"`
	Sort      int              `gorm:"column:sort;not null;default:0" json:"sort"`
	Status    int8             `gorm:"column:status;not null;default:1" json:"status"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (LogisticsCompany) TableName() string {
	return "logistics_companies"
}
