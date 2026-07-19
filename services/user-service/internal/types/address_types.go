package types

type AddressReq struct {
	ReceiverName  string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	Detail        string `json:"detail"`
	ProvinceCode  string `json:"province_code"`
	CityCode      string `json:"city_code"`
	DistrictCode  string `json:"district_code"`
	IsDefault     int    `json:"is_default"`
}
