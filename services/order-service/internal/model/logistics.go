package model

import "mymall/common"

type LogisticsCompany struct {
	ID        uint64           `db:"id" json:"id"`
	Name      string           `db:"name" json:"name"`
	Code      string           `db:"code" json:"code"`
	Sort      int              `db:"sort" json:"sort"`
	Status    int8             `db:"status" json:"status"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (LogisticsCompany) TableName() string {
	return "logistics_companies"
}
