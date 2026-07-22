package model

import "mymall/common"

type UserAddress struct {
	ID            uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	UserID        uint64           `gorm:"column:user_id;index" db:"user_id" json:"user_id"`
	ReceiverName  string           `gorm:"column:receiver_name;type:varchar(64)" db:"receiver_name" json:"receiver_name"`
	ReceiverPhone string           `gorm:"column:receiver_phone;type:varchar(20)" db:"receiver_phone" json:"receiver_phone"`
	Province      string           `gorm:"column:province;type:varchar(64)" db:"province" json:"province"`
	City          string           `gorm:"column:city;type:varchar(64)" db:"city" json:"city"`
	District      string           `gorm:"column:district;type:varchar(64)" db:"district" json:"district"`
	Detail        string           `gorm:"column:detail;type:varchar(255)" db:"detail" json:"detail"`
	ProvinceCode  string           `gorm:"column:province_code;type:varchar(12)" db:"province_code" json:"province_code"`
	CityCode      string           `gorm:"column:city_code;type:varchar(12)" db:"city_code" json:"city_code"`
	DistrictCode  string           `gorm:"column:district_code;type:varchar(12)" db:"district_code" json:"district_code"`
	IsDefault     int              `gorm:"column:is_default" db:"is_default" json:"is_default"`
	CreatedAt     common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt     common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (UserAddress) TableName() string { return "user_addresses" }

func (a *UserAddress) FullAddress() string {
	return a.Province + a.City + a.District + a.Detail
}
