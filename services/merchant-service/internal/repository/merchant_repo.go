package repository

import (
	"errors"

	"mymall/common/password"
	"mymall/services/merchant-service/internal/model"

	"gorm.io/gorm"
)

type MerchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) CreateApplication(app *model.ShopApplication) error {
	return r.db.Create(app).Error
}

func (r *MerchantRepository) FindPendingAppByUser(userID uint64) (*model.ShopApplication, error) {
	var app model.ShopApplication
	err := r.db.Where("user_id = ? AND status = ?", userID, model.AppPending).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *MerchantRepository) ListApplications(status string, page, pageSize int) ([]model.ShopApplication, int64, error) {
	q := r.db.Model(&model.ShopApplication{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.ShopApplication
	offset := (page - 1) * pageSize
	err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) FindApplication(id uint64) (*model.ShopApplication, error) {
	var app model.ShopApplication
	if err := r.db.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *MerchantRepository) ApproveApplication(appID, adminID uint64) (*model.Shop, error) {
	var shop *model.Shop
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var app model.ShopApplication
		if err := tx.First(&app, appID).Error; err != nil {
			return err
		}
		if app.Status != model.AppPending {
			return gorm.ErrInvalidData
		}
		shop = &model.Shop{
			Name:              app.ShopName,
			ContactName:       app.ContactName,
			ContactPhone:      app.ContactPhone,
			Description:       app.Description,
			Category:          app.Category,
			Province:          app.Province,
			City:              app.City,
			District:          app.District,
			Address:           app.Address,
			BusinessLicenseNo: app.BusinessLicenseNo,
			LegalPerson:       app.LegalPerson,
			LicenseImage:      app.LicenseImage,
			StorefrontImage:   app.StorefrontImage,
			OwnerUserID:       app.UserID,
			Status:            model.ShopApproved,
		}
		if err := tx.Create(shop).Error; err != nil {
			return err
		}
		member := model.ShopMember{
			ShopID:     shop.ID,
			UserID:     app.UserID,
			MemberRole: model.MemberOwner,
		}
		if err := tx.Create(&member).Error; err != nil {
			return err
		}
		if err := tx.Table("users").Where("id = ?", app.UserID).
			Update("role", "merchant_owner").Error; err != nil {
			return err
		}
		return tx.Model(&app).Updates(map[string]interface{}{
			"status":      model.AppApproved,
			"reviewed_by": adminID,
			"shop_id":     shop.ID,
			"reviewed_at": gorm.Expr("NOW()"),
		}).Error
	})
	return shop, err
}

func (r *MerchantRepository) RejectApplication(appID, adminID uint64, reason string) error {
	return r.db.Model(&model.ShopApplication{}).Where("id = ? AND status = ?", appID, model.AppPending).
		Updates(map[string]interface{}{
			"status":        model.AppRejected,
			"reject_reason": reason,
			"reviewed_by":   adminID,
			"reviewed_at":   gorm.Expr("NOW()"),
		}).Error
}

func (r *MerchantRepository) ListShops(status, name string, page, pageSize int) ([]model.Shop, int64, error) {
	q := r.db.Model(&model.Shop{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Shop
	offset := (page - 1) * pageSize
	err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) FindShop(id uint64) (*model.Shop, error) {
	var shop model.Shop
	if err := r.db.First(&shop, id).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *MerchantRepository) UpdateShopStatus(id uint64, status, reason string) error {
	updates := map[string]interface{}{"status": status}
	if reason != "" {
		updates["reject_reason"] = reason
	}
	return r.db.Model(&model.Shop{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MerchantRepository) UpdateShop(shop *model.Shop) error {
	return r.db.Model(shop).Updates(map[string]interface{}{
		"name":                shop.Name,
		"logo":                shop.Logo,
		"contact_name":        shop.ContactName,
		"contact_phone":       shop.ContactPhone,
		"description":         shop.Description,
		"category":            shop.Category,
		"province":            shop.Province,
		"city":                shop.City,
		"district":            shop.District,
		"address":             shop.Address,
		"business_license_no": shop.BusinessLicenseNo,
		"legal_person":        shop.LegalPerson,
		"license_image":       shop.LicenseImage,
		"storefront_image":    shop.StorefrontImage,
	}).Error
}

func (r *MerchantRepository) UpdateShopDisplay(shop *model.Shop) error {
	return r.db.Model(shop).Updates(map[string]interface{}{
		"name":             shop.Name,
		"logo":             shop.Logo,
		"contact_name":     shop.ContactName,
		"contact_phone":    shop.ContactPhone,
		"description":      shop.Description,
		"category":         shop.Category,
		"province":         shop.Province,
		"city":             shop.City,
		"district":         shop.District,
		"address":          shop.Address,
		"storefront_image": shop.StorefrontImage,
	}).Error
}

// CreateShopWithOwner 平台开店：绑定或新建店主账号
func (r *MerchantRepository) CreateShopWithOwner(shop *model.Shop, mobile, plainPwd, nickname string) (*model.Shop, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var ownerID uint64
		var existing struct {
			ID uint64 `gorm:"column:id"`
		}
		err := tx.Table("users").Select("id").Where("mobile = ?", mobile).First(&existing).Error
		if err == nil {
			ownerID = existing.ID
			if err := tx.Table("users").Where("id = ?", ownerID).Update("role", "merchant_owner").Error; err != nil {
				return err
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			if plainPwd == "" {
				return errors.New("新用户必须设置密码")
			}
			if nickname == "" {
				nickname = mobile
			}
			row := map[string]interface{}{
				"mobile":   mobile,
				"password": password.Hash(plainPwd),
				"nickname": nickname,
				"status":   1,
				"role":     "merchant_owner",
			}
			if err := tx.Table("users").Create(row).Error; err != nil {
				return err
			}
			if err := tx.Table("users").Select("id").Where("mobile = ?", mobile).First(&existing).Error; err != nil {
				return err
			}
			ownerID = existing.ID
		} else {
			return err
		}

		shop.OwnerUserID = ownerID
		shop.Status = model.ShopApproved
		if err := tx.Create(shop).Error; err != nil {
			return err
		}
		member := model.ShopMember{
			ShopID:     shop.ID,
			UserID:     ownerID,
			MemberRole: model.MemberOwner,
		}
		return tx.Create(&member).Error
	})
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (r *MerchantRepository) ResetOwnerPassword(shopID uint64, plainPwd string) error {
	shop, err := r.FindShop(shopID)
	if err != nil {
		return err
	}
	if shop.OwnerUserID == 0 {
		return errors.New("店铺无店主")
	}
	return r.db.Table("users").Where("id = ?", shop.OwnerUserID).
		Update("password", password.Hash(plainPwd)).Error
}

func (r *MerchantRepository) ListShopsByUser(userID uint64) ([]model.Shop, error) {
	var shops []model.Shop
	err := r.db.Table("shops").
		Joins("JOIN shop_members ON shop_members.shop_id = shops.id").
		Where("shop_members.user_id = ?", userID).
		Find(&shops).Error
	return shops, err
}

func (r *MerchantRepository) IsShopMember(shopID, userID uint64) bool {
	var count int64
	r.db.Model(&model.ShopMember{}).Where("shop_id = ? AND user_id = ?", shopID, userID).Count(&count)
	return count > 0
}
