package repository

import (
	"context"
	"errors"
	"time"

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

func (r *ProductRepository) GetList(ctx context.Context, page *pagination.PageReq) (map[string]interface{}, error) {
	return r.GetListByShop(ctx, page, 0, "on_sale")
}

func (r *ProductRepository) GetListByShop(ctx context.Context, page *pagination.PageReq, shopID uint64, status string) (map[string]interface{}, error) {
	return r.GetListFiltered(ctx, page, shopID, status, 0, "")
}

func (r *ProductRepository) categoryIDsIncludingChildren(ctx context.Context, categoryID uint64) []uint64 {
	if categoryID == 0 {
		return nil
	}
	ids := []uint64{categoryID}
	var children []uint64
	_ = r.db.WithContext(ctx).Model(&model.ProductCategory{}).Where("parent_id = ?", categoryID).Pluck("id", &children).Error
	ids = append(ids, children...)
	// 再下一层（最多三级）
	if len(children) > 0 {
		var grand []uint64
		_ = r.db.WithContext(ctx).Model(&model.ProductCategory{}).Where("parent_id IN ?", children).Pluck("id", &grand).Error
		ids = append(ids, grand...)
	}
	return ids
}

func (r *ProductRepository) GetListFiltered(ctx context.Context, page *pagination.PageReq, shopID uint64, status string, categoryID uint64, orderBy string) (map[string]interface{}, error) {
	empty := map[string]interface{}{
		"total": int64(0),
		"list":  []model.ProductListResp{},
	}

	apply := func(q *gorm.DB) *gorm.DB {
		if shopID > 0 {
			q = q.Where("shop_id = ?", shopID)
		}
		if status != "" {
			q = q.Where("status = ?", status)
		}
		if categoryID > 0 {
			ids := r.categoryIDsIncludingChildren(ctx, categoryID)
			q = q.Where("category_id IN ?", ids)
		}
		return q
	}

	q := apply(r.db.Model(&model.Product{}))
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return empty, err
	}
	if total == 0 {
		return empty, nil
	}

	query := apply(r.db.Model(&model.Product{}))
	switch orderBy {
	case "sold_count_desc":
		query = query.Order("sold_count DESC, id DESC")
	default:
		query = query.Order("id DESC")
	}
	result, err := pagination.Paginate[model.ProductListResp](query, page)
	if err != nil {
		return empty, err
	}
	return map[string]interface{}{
		"total":      result.Total,
		"page":       result.Page,
		"page_size":  result.PageSize,
		"total_page": result.TotalPage,
		"list":       result.List,
	}, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepository) UpdateByShop(ctx context.Context, id, shopID uint64, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ? AND shop_id = ?", id, shopID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("商品不存在或无权操作")
	}
	return nil
}

func (r *ProductRepository) ForceOffSale(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Update("status", "off_sale").Error
}

func (r *ProductRepository) GetDetail(ctx context.Context, id uint64) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) BatchGetByIDs(ctx context.Context, ids []uint64) ([]model.Product, error) {
	if len(ids) == 0 {
		return []model.Product{}, nil
	}
	var products []model.Product
	err := r.db.WithContext(ctx).Where("id IN ? AND status = ?", ids, "on_sale").Find(&products).Error
	return products, err
}

// ListSalesRank 今日销量优先，其次总销量；仅在售商品
func (r *ProductRepository) ListSalesRank(ctx context.Context, page, pageSize int) ([]model.ProductSalesRankItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "on_sale").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []model.ProductSalesRankItem{}, 0, nil
	}

	var list []model.ProductSalesRankItem
	err := r.db.WithContext(ctx).Raw(`
SELECT p.id, p.shop_id, p.name, COALESCE(s.name,'') AS shop_name,
       COALESCE(p.main_image,'') AS main_image, p.sale_price,
       p.sold_count, COALESCE(t.today_sold, 0) AS today_sold
FROM products p
LEFT JOIN shops s ON s.id = p.shop_id
LEFT JOIN (
  SELECT oi.product_id, SUM(oi.quantity) AS today_sold
  FROM order_items oi
  INNER JOIN orders o ON o.id = oi.order_id
  WHERE o.status IN ('confirmed','shipped','completed','reviewed')
    AND o.created_at >= ? AND o.created_at < ?
  GROUP BY oi.product_id
) t ON t.product_id = p.id
WHERE p.status = 'on_sale'
ORDER BY COALESCE(t.today_sold, 0) DESC, p.sold_count DESC, p.id DESC
LIMIT ? OFFSET ?`, dayStart, dayEnd, pageSize, (page-1)*pageSize).Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *ProductRepository) ReserveStock(ctx context.Context, items []StockItem) error {
	admin := NewProductAdminRepository(r.db.WithContext(ctx))
	skuItems := make([]SkuStockItem, 0, len(items))
	for _, it := range items {
		skuItems = append(skuItems, SkuStockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
	}
	return admin.ReserveSkuStock(ctx, skuItems)
}

func (r *ProductRepository) ReleaseStock(ctx context.Context, items []StockItem) error {
	admin := NewProductAdminRepository(r.db.WithContext(ctx))
	skuItems := make([]SkuStockItem, 0, len(items))
	for _, it := range items {
		skuItems = append(skuItems, SkuStockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
	}
	return admin.ReleaseSkuStock(ctx, skuItems)
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetList(ctx context.Context, page *pagination.PageReq) (*pagination.PageRes[model.ProductCategory], error) {
	res := &pagination.PageRes[model.ProductCategory]{List: []model.ProductCategory{}}

	var total int64
	if err := r.db.WithContext(ctx).Model(&model.ProductCategory{}).Where("is_show = ?", true).Count(&total).Error; err != nil {
		return res, err
	}
	if total == 0 {
		return res, nil
	}

	query := r.db.WithContext(ctx).Model(&model.ProductCategory{}).Where("is_show = ?", true).Order("sort_order ASC, id ASC")
	return pagination.Paginate[model.ProductCategory](query, page)
}

// ListAll 管理端：含隐藏分类，按排序返回全量
func (r *CategoryRepository) ListAll(ctx context.Context) ([]model.ProductCategory, error) {
	var list []model.ProductCategory
	err := r.db.WithContext(ctx).Model(&model.ProductCategory{}).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *CategoryRepository) GetDetail(ctx context.Context, id uint64) (*model.ProductCategory, error) {
	var category model.ProductCategory
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *model.ProductCategory) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CategoryRepository) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.ProductCategory{}).Where("id = ?", id).Updates(updates).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint64) error {
	var child int64
	_ = r.db.WithContext(ctx).Model(&model.ProductCategory{}).Where("parent_id = ?", id).Count(&child).Error
	if child > 0 {
		return errors.New("请先删除子分类")
	}
	var used int64
	_ = r.db.WithContext(ctx).Model(&model.Product{}).Where("category_id = ? AND status <> ?", id, "deleted").Count(&used).Error
	if used > 0 {
		return errors.New("分类下仍有商品，无法删除")
	}
	return r.db.WithContext(ctx).Delete(&model.ProductCategory{}, id).Error
}
