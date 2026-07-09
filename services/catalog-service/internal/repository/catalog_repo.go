package repository

import (
	"mymall/pkg/pagination"
	"mymall/services/catalog-service/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetList(page *pagination.PageReq) (map[string]interface{}, error) {
	res := map[string]interface{}{
		"total": int64(0),
		"data":  []model.ProductListResp{},
	}

	var total int64
	if err := r.db.Model(&model.Product{}).Count(&total).Error; err != nil {
		return res, err
	}
	if total == 0 {
		return res, nil
	}

	query := r.db.Model(&model.Product{}).Where("status = ?", "on_sale")
	result, err := pagination.Paginate[model.ProductListResp](query, page)
	if err != nil {
		return res, err
	}
	res["data"] = result
	return res, nil
}

func (r *ProductRepository) GetDetail(id uint64) (*model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetList(page *pagination.PageReq) (*pagination.PageRes[model.ProductCategory], error) {
	res := &pagination.PageRes[model.ProductCategory]{List: []model.ProductCategory{}}

	var total int64
	if err := r.db.Model(&model.ProductCategory{}).Count(&total).Error; err != nil {
		return res, err
	}
	if total == 0 {
		return res, nil
	}

	query := r.db.Model(&model.ProductCategory{}).Where("is_show = ?", true)
	return pagination.Paginate[model.ProductCategory](query, page)
}

func (r *CategoryRepository) GetDetail(id uint64) (*model.ProductCategory, error) {
	var category model.ProductCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
