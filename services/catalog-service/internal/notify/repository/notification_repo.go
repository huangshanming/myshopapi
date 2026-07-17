package repository

import (
	"errors"

	"mymall/services/catalog-service/internal/notify/model"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *model.ShopNotification) error {
	return r.db.Create(n).Error
}

type NotificationListFilter struct {
	ShopID   uint64
	IsRead   *int8 // nil=全部
	Page     int
	PageSize int
}

func (r *NotificationRepository) List(f NotificationListFilter) ([]model.ShopNotification, int64, error) {
	q := r.db.Model(&model.ShopNotification{}).Where("shop_id = ?", f.ShopID)
	if f.IsRead != nil {
		q = q.Where("is_read = ?", *f.IsRead)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	var list []model.ShopNotification
	err := q.Order("id DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&list).Error
	return list, total, err
}

func (r *NotificationRepository) UnreadCount(shopID uint64) (int64, error) {
	var cnt int64
	err := r.db.Model(&model.ShopNotification{}).
		Where("shop_id = ? AND is_read = 0", shopID).Count(&cnt).Error
	return cnt, err
}

func (r *NotificationRepository) MarkRead(id, shopID uint64) error {
	res := r.db.Model(&model.ShopNotification{}).
		Where("id = ? AND shop_id = ?", id, shopID).
		Update("is_read", 1)
	if res.RowsAffected == 0 {
		return errors.New("通知不存在")
	}
	return res.Error
}

func (r *NotificationRepository) MarkAllRead(shopID uint64) error {
	return r.db.Model(&model.ShopNotification{}).
		Where("shop_id = ? AND is_read = 0", shopID).
		Update("is_read", 1).Error
}
