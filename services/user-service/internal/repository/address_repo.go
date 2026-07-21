package repository

import (
	"context"
	"errors"

	"mymall/services/user-service/internal/model"

	"gorm.io/gorm"
)

func (r *UserRepository) ListAddresses(ctx context.Context, userID uint64) ([]model.UserAddress, error) {
	var list []model.UserAddress
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, id DESC").Find(&list).Error
	return list, err
}

func (r *UserRepository) GetAddress(ctx context.Context, userID, id uint64) (*model.UserAddress, error) {
	var a model.UserAddress
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *UserRepository) CreateAddress(ctx context.Context, a *model.UserAddress) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if a.IsDefault == 1 {
			if err := tx.Model(&model.UserAddress{}).Where("user_id = ?", a.UserID).Update("is_default", 0).Error; err != nil {
				return err
			}
		} else {
			var n int64
			if err := tx.Model(&model.UserAddress{}).Where("user_id = ?", a.UserID).Count(&n).Error; err != nil {
				return err
			}
			if n == 0 {
				a.IsDefault = 1
			}
		}
		return tx.Create(a).Error
	})
}

func (r *UserRepository) UpdateAddress(ctx context.Context, userID, id uint64, a *model.UserAddress) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.UserAddress
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&existing).Error; err != nil {
			return err
		}
		if a.IsDefault == 1 {
			if err := tx.Model(&model.UserAddress{}).Where("user_id = ?", userID).Update("is_default", 0).Error; err != nil {
				return err
			}
		}
		return tx.Model(&existing).Updates(map[string]interface{}{
			"receiver_name":  a.ReceiverName,
			"receiver_phone": a.ReceiverPhone,
			"province":       a.Province,
			"city":           a.City,
			"district":       a.District,
			"detail":         a.Detail,
			"province_code":  a.ProvinceCode,
			"city_code":      a.CityCode,
			"district_code":  a.DistrictCode,
			"is_default":     a.IsDefault,
		}).Error
	})
}

func (r *UserRepository) DeleteAddress(ctx context.Context, userID, id uint64) error {
	res := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.UserAddress{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepository) SetDefaultAddress(ctx context.Context, userID, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.UserAddress
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&existing).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.UserAddress{}).Where("user_id = ?", userID).Update("is_default", 0).Error; err != nil {
			return err
		}
		return tx.Model(&existing).Update("is_default", 1).Error
	})
}

func (r *UserRepository) GetAddressByID(ctx context.Context, userID, id uint64) (*model.UserAddress, error) {
	a, err := r.GetAddress(ctx, userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("收货地址不存在")
		}
		return nil, err
	}
	return a, nil
}
