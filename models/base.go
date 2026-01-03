package model

import "mymall/common"

type BaseModel struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement;comment:ID" json:"id"`
	CreatedAt common.LocalTime `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}
