package model

import "mymall/common"

// UserAddress runtime entity (json tags for mall-uni; not modelgen alias).
type UserAddress struct {
	ID            uint64           `db:"id" json:"id"`
	UserID        uint64           `db:"user_id" json:"user_id"`
	ReceiverName  string           `db:"receiver_name" json:"receiver_name"`
	ReceiverPhone string           `db:"receiver_phone" json:"receiver_phone"`
	Province      string           `db:"province" json:"province"`
	City          string           `db:"city" json:"city"`
	District      string           `db:"district" json:"district"`
	Detail        string           `db:"detail" json:"detail"`
	ProvinceCode  string           `db:"province_code" json:"province_code"`
	CityCode      string           `db:"city_code" json:"city_code"`
	DistrictCode  string           `db:"district_code" json:"district_code"`
	IsDefault     int              `db:"is_default" json:"is_default"`
	CreatedAt     common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt     common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (a *UserAddress) FullAddress() string {
	if a == nil {
		return ""
	}
	return a.Province + a.City + a.District + a.Detail
}
