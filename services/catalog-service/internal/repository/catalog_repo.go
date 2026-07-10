package repository

import (
	"errors"

	"mymall/pkg/pagination"
	"mymall/services/catalog-service/internal/model"

	"gorm.io/gorm"
)

type StockItem struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

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

func (r *ProductRepository) BatchGetByIDs(ids []uint64) ([]model.Product, error) {
	if len(ids) == 0 {
		return []model.Product{}, nil
	}
	var products []model.Product
	err := r.db.Where("id IN ? AND status = ?", ids, "on_sale").Find(&products).Error
	return products, err
}

func (r *ProductRepository) ReserveStock(items []StockItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			result := tx.Model(&model.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				Update("stock", gorm.Expr("stock - ?", item.Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("insufficient stock")
			}
		}
		return nil
	})
}

func (r *ProductRepository) ReleaseStock(items []StockItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}
		return nil
	})
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
