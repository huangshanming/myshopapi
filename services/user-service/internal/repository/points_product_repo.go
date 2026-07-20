package repository

import (
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

func (r *PointsProductRepository) List(page, pageSize int, status, keyword string) ([]model.PointsProduct, int64, error) {
	q := r.db.Model(&model.PointsProduct{})
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

func (r *PointsProductRepository) GetByID(id uint64) (*model.PointsProduct, error) {
	var p model.PointsProduct
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PointsProductRepository) Create(p *model.PointsProduct) error {
	return r.db.Create(p).Error
}

func (r *PointsProductRepository) Update(id uint64, updates map[string]interface{}) error {
	return r.db.Model(&model.PointsProduct{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PointsProductRepository) Delete(id uint64) error {
	return r.db.Delete(&model.PointsProduct{}, id).Error
}
