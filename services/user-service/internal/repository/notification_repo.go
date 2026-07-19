package repository

import (
	"errors"

	"mymall/services/user-service/internal/model"
)

func (r *UserRepository) CreateNotification(n *model.UserNotification) error {
	return r.db.Create(n).Error
}

func (r *UserRepository) CreateNotifications(list []model.UserNotification) error {
	if len(list) == 0 {
		return nil
	}
	return r.db.CreateInBatches(list, 200).Error
}

func (r *UserRepository) ListNotifications(userID uint64, page, pageSize int) ([]model.UserNotification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.Model(&model.UserNotification{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.UserNotification
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *UserRepository) UnreadNotificationCount(userID uint64) (int64, error) {
	var n int64
	err := r.db.Model(&model.UserNotification{}).Where("user_id = ? AND is_read = 0", userID).Count(&n).Error
	return n, err
}

func (r *UserRepository) MarkNotificationRead(userID, id uint64) error {
	res := r.db.Model(&model.UserNotification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", 1)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("消息不存在")
	}
	return nil
}

func (r *UserRepository) MarkAllNotificationsRead(userID uint64) error {
	return r.db.Model(&model.UserNotification{}).
		Where("user_id = ? AND is_read = 0", userID).
		Update("is_read", 1).Error
}

func (r *UserRepository) ListActiveUserIDs(offset, limit int) ([]uint64, error) {
	var ids []uint64
	err := r.db.Model(&model.User{}).Where("status = 1").Order("id ASC").
		Offset(offset).Limit(limit).Pluck("id", &ids).Error
	return ids, err
}

func (r *UserRepository) CreateNotificationBatch(b *model.UserNotificationBatch) error {
	return r.db.Create(b).Error
}

func (r *UserRepository) UpdateNotificationBatchSuccess(id uint64, success int) error {
	return r.db.Model(&model.UserNotificationBatch{}).Where("id = ?", id).Update("success_count", success).Error
}

func (r *UserRepository) ListNotificationBatches(page, pageSize int) ([]model.UserNotificationBatch, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.Model(&model.UserNotificationBatch{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.UserNotificationBatch
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *UserRepository) GetNotificationBatch(id uint64) (*model.UserNotificationBatch, error) {
	var b model.UserNotificationBatch
	if err := r.db.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

type BatchRecipientRow struct {
	UserID   uint64 `json:"user_id"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
	IsRead   int8   `json:"is_read"`
}

func (r *UserRepository) ListBatchRecipients(batchID uint64, page, pageSize int) ([]BatchRecipientRow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.Table("user_notifications AS n").Where("n.batch_id = ?", batchID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []BatchRecipientRow
	err := r.db.Table("user_notifications AS n").
		Select("n.user_id, COALESCE(u.nickname,'') AS nickname, COALESCE(u.mobile,'') AS mobile, n.is_read").
		Joins("LEFT JOIN users u ON u.id = n.user_id").
		Where("n.batch_id = ?", batchID).
		Order("n.id ASC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Scan(&list).Error
	return list, total, err
}
