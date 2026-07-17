package repository

import (
	"errors"

	"mymall/pkg/pagination"
	"mymall/services/catalog-service/internal/product/model"

	"gorm.io/gorm"
)

type StockItem struct {
	ProductID uint64 `json:"product_id"`
	SkuID     uint64 `json:"sku_id"`
	Quantity  int    `json:"quantity"`
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetList(page *pagination.PageReq) (map[string]interface{}, error) {
	return r.GetListByShop(page, 0, "on_sale")
}

func (r *ProductRepository) GetListByShop(page *pagination.PageReq, shopID uint64, status string) (map[string]interface{}, error) {
	empty := map[string]interface{}{
		"total": int64(0),
		"list":  []model.ProductListResp{},
	}

	q := r.db.Model(&model.Product{})
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return empty, err
	}
	if total == 0 {
		return empty, nil
	}

	query := r.db.Model(&model.Product{})
	if shopID > 0 {
		query = query.Where("shop_id = ?", shopID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	result, err := pagination.Paginate[model.ProductListResp](query, page)
	if err != nil {
		return empty, err
	}
	// 统一 { total, list }，避免再包一层 data 导致前端 el-table 拿到对象
	return map[string]interface{}{
		"total":      result.Total,
		"page":       result.Page,
		"page_size":  result.PageSize,
		"total_page": result.TotalPage,
		"list":       result.List,
	}, nil
}

func (r *ProductRepository) Create(p *model.Product) error {
	return r.db.Create(p).Error
}

func (r *ProductRepository) UpdateByShop(id, shopID uint64, updates map[string]interface{}) error {
	result := r.db.Model(&model.Product{}).Where("id = ? AND shop_id = ?", id, shopID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("商品不存在或无权操作")
	}
	return nil
}

func (r *ProductRepository) ForceOffSale(id uint64) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Update("status", "off_sale").Error
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
	admin := NewProductAdminRepository(r.db)
	skuItems := make([]SkuStockItem, 0, len(items))
	for _, it := range items {
		skuItems = append(skuItems, SkuStockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
	}
	return admin.ReserveSkuStock(skuItems)
}

func (r *ProductRepository) ReleaseStock(items []StockItem) error {
	admin := NewProductAdminRepository(r.db)
	skuItems := make([]SkuStockItem, 0, len(items))
	for _, it := range items {
		skuItems = append(skuItems, SkuStockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
	}
	return admin.ReleaseSkuStock(skuItems)
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
	if err := r.db.Model(&model.ProductCategory{}).Where("is_show = ?", true).Count(&total).Error; err != nil {
		return res, err
	}
	if total == 0 {
		return res, nil
	}

	query := r.db.Model(&model.ProductCategory{}).Where("is_show = ?", true)
	return pagination.Paginate[model.ProductCategory](query, page)
}

// ListAll 管理端：含隐藏分类，按排序返回全量
func (r *CategoryRepository) ListAll() ([]model.ProductCategory, error) {
	var list []model.ProductCategory
	err := r.db.Model(&model.ProductCategory{}).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *CategoryRepository) GetDetail(id uint64) (*model.ProductCategory, error) {
	var category model.ProductCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(c *model.ProductCategory) error {
	return r.db.Create(c).Error
}

func (r *CategoryRepository) Update(id uint64, updates map[string]interface{}) error {
	return r.db.Model(&model.ProductCategory{}).Where("id = ?", id).Updates(updates).Error
}

func (r *CategoryRepository) Delete(id uint64) error {
	var child int64
	_ = r.db.Model(&model.ProductCategory{}).Where("parent_id = ?", id).Count(&child).Error
	if child > 0 {
		return errors.New("请先删除子分类")
	}
	var used int64
	_ = r.db.Model(&model.Product{}).Where("category_id = ? AND status <> ?", id, "deleted").Count(&used).Error
	if used > 0 {
		return errors.New("分类下仍有商品，无法删除")
	}
	return r.db.Delete(&model.ProductCategory{}, id).Error
}
