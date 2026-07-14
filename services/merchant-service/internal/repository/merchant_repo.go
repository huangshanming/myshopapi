package repository

import (
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
			Name:         app.ShopName,
			ContactName:  app.ContactName,
			ContactPhone: app.ContactPhone,
			Description:  app.Description,
			OwnerUserID:  app.UserID,
			Status:       model.ShopApproved,
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
		// 提升用户角色
		if err := tx.Table("users").Where("id = ?", app.UserID).
			Update("role", "merchant_owner").Error; err != nil {
			return err
		}
		now := app.UpdatedAt
		_ = now
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

func (r *MerchantRepository) ListShops(status string, page, pageSize int) ([]model.Shop, int64, error) {
	q := r.db.Model(&model.Shop{})
	if status != "" {
		q = q.Where("status = ?", status)
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
		"name":          shop.Name,
		"logo":          shop.Logo,
		"contact_name":  shop.ContactName,
		"contact_phone": shop.ContactPhone,
		"description":   shop.Description,
	}).Error
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
