package model

import "mymall/services/user-service/internal/modelgen"

// UserAddress is a defined type so FullAddress() can be attached.
type UserAddress modelgen.UserAddress

func (a *UserAddress) FullAddress() string {
	if a == nil {
		return ""
	}
	return a.Province + a.City + a.District + a.Detail
}
