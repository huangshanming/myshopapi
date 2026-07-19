package model

import "mymall/common"

type UserAddress struct {
	ID            uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        uint64           `gorm:"column:user_id;index" json:"user_id"`
	ReceiverName  string           `gorm:"column:receiver_name;type:varchar(64)" json:"receiver_name"`
	ReceiverPhone string           `gorm:"column:receiver_phone;type:varchar(20)" json:"receiver_phone"`
	Province      string           `gorm:"column:province;type:varchar(64)" json:"province"`
	City          string           `gorm:"column:city;type:varchar(64)" json:"city"`
	District      string           `gorm:"column:district;type:varchar(64)" json:"district"`
	Detail        string           `gorm:"column:detail;type:varchar(255)" json:"detail"`
	IsDefault     int              `gorm:"column:is_default" json:"is_default"`
	CreatedAt     common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (UserAddress) TableName() string { return "user_addresses" }

func (a *UserAddress) FullAddress() string {
	return a.Province + a.City + a.District + a.Detail
}
