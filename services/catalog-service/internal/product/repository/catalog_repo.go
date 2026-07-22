package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"mymall/pkg/pagination"
	"mymall/services/catalog-service/internal/product/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	// copy_from_id / numeric nullables use IFNULL: go-zero sqlx pre-allocates *T so NULL cannot scan into *uint64.
	productColumns = "id, shop_id, product_no, name, IFNULL(subtitle,'') AS subtitle, IFNULL(description,'') AS description, IFNULL(main_image,'') AS main_image, image_list, IFNULL(video_url,'') AS video_url, IFNULL(market_price,0) AS market_price, sale_price, IFNULL(cost_price,0) AS cost_price, discount, IFNULL(discount_price,0) AS discount_price, stock, stock_warn, sold_count, view_count, collect_count, avg_rating, review_count, good_rate, pet_type, pet_age, pet_size, IFNULL(weight,0) AS weight, IFNULL(unit,'') AS unit, IFNULL(brand_id,0) AS brand_id, category_id, product_type, IFNULL(spec_json,'') AS spec_json, tags, nutrition_info, IFNULL(ingredients,'') AS ingredients, IFNULL(feeding_guide,'') AS feeding_guide, IFNULL(shelf_life,0) AS shelf_life, IFNULL(storage_condition,'') AS storage_condition, status, is_hot, is_new, is_recommend, is_prescription, is_imported, is_organic, is_grain_free, publish_time, schedule_on_at, schedule_off_at, IFNULL(copy_from_id,0) AS copy_from_id, deleted_at, created_at, updated_at"
	productListColumns = "id, shop_id, product_no, name, IFNULL(subtitle,'') AS subtitle, IFNULL(main_image,'') AS main_image, IFNULL(market_price,0) AS market_price, sale_price, discount, stock, sold_count, collect_count, avg_rating, review_count, good_rate, category_id, IFNULL(brand_id,0) AS brand_id, pet_type, tags, status, is_hot, is_new, is_recommend, is_prescription, publish_time, created_at"
	productCategoryColumns = "id, parent_id, name, IFNULL(icon,'') AS icon, IFNULL(image,'') AS image, IFNULL(description,'') AS description, sort_order, level, is_show, product_count, created_at, updated_at"
)

type StockItem struct {
	ProductID uint64 `json:"product_id"`
	SkuID     uint64 `json:"sku_id"`
	Quantity  int    `json:"quantity"`
}

type ProductRepository struct {
	conn sqlx.SqlConn
}

func NewProductRepository(conn sqlx.SqlConn) *ProductRepository {
	return &ProductRepository{conn: conn}
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
	_ = r.conn.QueryRowsPartialCtx(ctx, &children, "SELECT id FROM product_categories WHERE parent_id=?", categoryID)
	ids = append(ids, children...)
	if len(children) > 0 {
		var grand []uint64
		_ = r.conn.QueryRowsPartialCtx(ctx, &grand,
			"SELECT id FROM product_categories WHERE parent_id IN ("+placeholders(len(children))+")",
			inArgs(children)...,
		)
		ids = append(ids, grand...)
	}
	return ids
}

func (r *ProductRepository) buildListWhere(ctx context.Context, shopID uint64, status string, categoryID uint64) (string, []any) {
	where := []string{"1=1"}
	args := make([]any, 0, 4)
	if shopID > 0 {
		where = append(where, "shop_id=?")
		args = append(args, shopID)
	}
	if status != "" {
		where = append(where, "status=?")
		args = append(args, status)
	}
	if categoryID > 0 {
		ids := r.categoryIDsIncludingChildren(ctx, categoryID)
		where = append(where, "category_id IN ("+placeholders(len(ids))+")")
		args = append(args, inArgs(ids)...)
	}
	return strings.Join(where, " AND "), args
}

func (r *ProductRepository) GetListFiltered(ctx context.Context, page *pagination.PageReq, shopID uint64, status string, categoryID uint64, orderBy string) (map[string]interface{}, error) {
	empty := map[string]interface{}{
		"total": int64(0),
		"list":  []model.ProductListResp{},
	}
	whereSQL, args := r.buildListWhere(ctx, shopID, status, categoryID)

	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM products WHERE "+whereSQL, args...)
	if err != nil {
		return empty, err
	}
	if total == 0 {
		return empty, nil
	}

	pg, pageSize, offset := pagination.Normalize(page)
	order := "id DESC"
	if orderBy == "sold_count_desc" {
		order = "sold_count DESC, id DESC"
	}

	var list []model.ProductListResp
	qArgs := append(args, pageSize, offset)
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+productListColumns+" FROM products WHERE "+whereSQL+" ORDER BY "+order+" LIMIT ? OFFSET ?",
		qArgs...,
	)
	if err != nil {
		return empty, err
	}
	totalPage := (int(total) + pageSize - 1) / pageSize
	return map[string]interface{}{
		"total":      total,
		"page":       pg,
		"page_size":  pageSize,
		"total_page": totalPage,
		"list":       list,
	}, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) error {
	id, err := lastInsertID(ctx, r.conn,
		`INSERT INTO products (shop_id, product_no, name, subtitle, main_image, sale_price, stock, category_id, pet_type, discount, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ShopID, p.ProductNo, p.Name, p.Subtitle, p.MainImage, p.SalePrice, p.Stock, p.CategoryID, p.PetType, p.Discount, p.Status,
	)
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func (r *ProductRepository) UpdateByShop(ctx context.Context, id, shopID uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("products", updates, "id=? AND shop_id=?", id, shopID)
	if err != nil {
		return err
	}
	n, err := execAffected(ctx, r.conn, query, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("商品不存在或无权操作")
	}
	return nil
}

func (r *ProductRepository) ForceOffSale(ctx context.Context, id uint64) error {
	_, err := r.conn.ExecCtx(ctx, "UPDATE products SET status='off_sale' WHERE id=?", id)
	return err
}

func (r *ProductRepository) GetDetail(ctx context.Context, id uint64) (*model.Product, error) {
	var product model.Product
	err := r.conn.QueryRowPartialCtx(ctx, &product, "SELECT "+productColumns+" FROM products WHERE id=? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetCollectCount(ctx context.Context, productID uint64) (int64, error) {
	var n int64
	err := r.conn.QueryRowPartialCtx(ctx, &n, "SELECT collect_count FROM products WHERE id=? LIMIT 1", productID)
	return int64(n), err
}

func (r *ProductRepository) BatchGetByIDs(ctx context.Context, ids []uint64) ([]model.Product, error) {
	if len(ids) == 0 {
		return []model.Product{}, nil
	}
	var products []model.Product
	err := r.conn.QueryRowsPartialCtx(ctx, &products,
		"SELECT "+productColumns+" FROM products WHERE id IN ("+placeholders(len(ids))+") AND status='on_sale'",
		inArgs(ids)...,
	)
	return products, err
}

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

	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM products WHERE status='on_sale'")
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []model.ProductSalesRankItem{}, 0, nil
	}

	var list []model.ProductSalesRankItem
	err = r.conn.QueryRowsPartialCtx(ctx, &list, `
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
LIMIT ? OFFSET ?`, dayStart, dayEnd, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *ProductRepository) ReserveStock(ctx context.Context, items []StockItem) error {
	admin := NewProductAdminRepository(r.conn)
	skuItems := make([]SkuStockItem, 0, len(items))
	for _, it := range items {
		skuItems = append(skuItems, SkuStockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
	}
	return admin.ReserveSkuStock(ctx, skuItems)
}

func (r *ProductRepository) ReleaseStock(ctx context.Context, items []StockItem) error {
	admin := NewProductAdminRepository(r.conn)
	skuItems := make([]SkuStockItem, 0, len(items))
	for _, it := range items {
		skuItems = append(skuItems, SkuStockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
	}
	return admin.ReleaseSkuStock(ctx, skuItems)
}

type CategoryRepository struct {
	conn sqlx.SqlConn
}

func NewCategoryRepository(conn sqlx.SqlConn) *CategoryRepository {
	return &CategoryRepository{conn: conn}
}

func (r *CategoryRepository) GetList(ctx context.Context, page *pagination.PageReq) (*pagination.PageRes[model.ProductCategory], error) {
	res := &pagination.PageRes[model.ProductCategory]{List: []model.ProductCategory{}}

	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM product_categories WHERE is_show=1")
	if err != nil {
		return res, err
	}
	if total == 0 {
		return res, nil
	}

	pg, pageSize, offset := pagination.Normalize(page)
	var list []model.ProductCategory
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+productCategoryColumns+" FROM product_categories WHERE is_show=1 ORDER BY sort_order ASC, id ASC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return res, err
	}
	totalPage := (int(total) + pageSize - 1) / pageSize
	res.Total = total
	res.Page = pg
	res.PageSize = pageSize
	res.TotalPage = totalPage
	res.List = list
	return res, nil
}

func (r *CategoryRepository) ListAll(ctx context.Context) ([]model.ProductCategory, error) {
	var list []model.ProductCategory
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+productCategoryColumns+" FROM product_categories ORDER BY sort_order ASC, id ASC")
	return list, err
}

func (r *CategoryRepository) GetDetail(ctx context.Context, id uint64) (*model.ProductCategory, error) {
	var category model.ProductCategory
	err := r.conn.QueryRowPartialCtx(ctx, &category, "SELECT "+productCategoryColumns+" FROM product_categories WHERE id=? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *model.ProductCategory) error {
	id, err := lastInsertID(ctx, r.conn,
		`INSERT INTO product_categories (parent_id, name, icon, image, description, sort_order, level, is_show, product_count)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ParentId, c.Name, c.Icon, c.Image, c.Description, c.SortOrder, c.Level, c.IsShow, c.ProductCount,
	)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *CategoryRepository) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("product_categories", updates, "id=?", id)
	if err != nil {
		return err
	}
	_, err = r.conn.ExecCtx(ctx, query, args...)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint64) error {
	child, _ := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM product_categories WHERE parent_id=?", id)
	if child > 0 {
		return errors.New("请先删除子分类")
	}
	used, _ := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM products WHERE category_id=? AND status<>?", id, "deleted")
	if used > 0 {
		return errors.New("分类下仍有商品，无法删除")
	}
	_, err := r.conn.ExecCtx(ctx, "DELETE FROM product_categories WHERE id=?", id)
	return err
}
