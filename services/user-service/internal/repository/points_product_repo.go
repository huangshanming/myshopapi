package repository

import (
	"context"
	"strings"

	"mymall/services/user-service/internal/model"

	"gorm.io/gorm"
)

type PointsProductRepository struct {
	db *gorm.DB
}

func NewPointsProductRepository(db *gorm.DB) *PointsProductRepository {
	return &PointsProductRepository{db: db}
}

func (r *PointsProductRepository) List(ctx context.Context, page, pageSize int, status, keyword string) ([]model.PointsProduct, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.PointsProduct{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if kw := strings.TrimSpace(keyword); kw != "" {
		q = q.Where("name LIKE ?", "%"+kw+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.PointsProduct
	err := q.Order("sort DESC, id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error
	return list, total, err
}

func (r *PointsProductRepository) GetByID(ctx context.Context, id uint64) (*model.PointsProduct, error) {
	var p model.PointsProduct
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PointsProductRepository) Create(ctx context.Context, p *model.PointsProduct) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PointsProductRepository) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.PointsProduct{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PointsProductRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.PointsProduct{}, id).Error
}
