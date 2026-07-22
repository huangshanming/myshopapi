package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/product/skuutil"
	"mymall/services/catalog-service/internal/product/types"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	productSkuColumns = "id, product_id, shop_id, sku_no, spec_values, spec_key, sale_price, market_price, cost_price, stock, stock_warn, barcode, status, sold_count, created_at, updated_at, deleted_at"
	productImageColumns = "id, product_id, shop_id, url, typ, sort, created_at, updated_at"
	productAttrColumns = "id, product_id, template_id, attr_key, attr_label, attr_value"
	productTagColumns = "id, shop_id, name, color, status, created_at, updated_at"
	productAttrTemplateColumns = "id, shop_id, name, attrs_json, status, created_at, updated_at"
	productScheduleColumns = "id, product_id, shop_id, action, run_at, status, locked_at, created_at, updated_at"
	productBatchJobColumns = "id, shop_id, job_type, payload_json, progress, total, status, result_msg, operator_id, created_at, updated_at"
	productOpLogColumns = "id, shop_id, product_id, operator_id, action, before_json, after_json, ip, created_at"
)

type ProductAdminRepository struct {
	conn sqlx.SqlConn
}

func NewProductAdminRepository(conn sqlx.SqlConn) *ProductAdminRepository {
	return &ProductAdminRepository{conn: conn}
}

type ProductListFilter struct {
	ShopID        uint64
	Name          string
	ProductNo     string
	CategoryID    uint64
	Status        string
	ProductType   string
	StockWarnOnly bool
	Page          int
	PageSize      int
	OrderBy       string
	Recycle       bool
	CreatedFrom   *time.Time
	CreatedTo     *time.Time
	PublishFrom   *time.Time
	PublishTo     *time.Time
	PlatformScope bool
}

func (r *ProductAdminRepository) List(ctx context.Context, f ProductListFilter) ([]model.Product, int64, error) {
	where := []string{"1=1"}
	args := make([]any, 0, 16)
	if f.PlatformScope {
		if f.ShopID > 0 {
			where = append(where, "shop_id=?")
			args = append(args, f.ShopID)
		}
	} else {
		where = append(where, "shop_id=?")
		args = append(args, f.ShopID)
	}
	if f.Recycle {
		where = append(where, "status=?")
		args = append(args, model.ProductDeleted)
	} else {
		where = append(where, "status<>?")
		args = append(args, model.ProductDeleted)
	}
	if f.Name != "" {
		where = append(where, "name LIKE ?")
		args = append(args, "%"+f.Name+"%")
	}
	if f.ProductNo != "" {
		where = append(where, "product_no LIKE ?")
		args = append(args, "%"+f.ProductNo+"%")
	}
	if f.CategoryID > 0 {
		where = append(where, "category_id=?")
		args = append(args, f.CategoryID)
	}
	if f.Status != "" {
		where = append(where, "status=?")
		args = append(args, f.Status)
	}
	if f.ProductType != "" {
		where = append(where, "product_type=?")
		args = append(args, f.ProductType)
	}
	if f.StockWarnOnly {
		where = append(where, "stock <= stock_warn")
	}
	if f.CreatedFrom != nil {
		where = append(where, "created_at >= ?")
		args = append(args, *f.CreatedFrom)
	}
	if f.CreatedTo != nil {
		where = append(where, "created_at <= ?")
		args = append(args, *f.CreatedTo)
	}
	if f.PublishFrom != nil {
		where = append(where, "publish_time >= ?")
		args = append(args, *f.PublishFrom)
	}
	if f.PublishTo != nil {
		where = append(where, "publish_time <= ?")
		args = append(args, *f.PublishTo)
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM products WHERE "+whereSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	order := "id DESC"
	switch f.OrderBy {
	case "sale_price_asc":
		order = "sale_price ASC"
	case "sale_price_desc":
		order = "sale_price DESC"
	case "stock_asc":
		order = "stock ASC"
	case "sold_desc":
		order = "sold_count DESC"
	case "collect_desc":
		order = "collect_count DESC"
	case "collect_asc":
		order = "collect_count ASC"
	case "created_asc":
		order = "id ASC"
	}
	var list []model.Product
	qArgs := append(args, f.PageSize, (f.Page-1)*f.PageSize)
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+productColumns+" FROM products WHERE "+whereSQL+" ORDER BY "+order+" LIMIT ? OFFSET ?",
		qArgs...,
	)
	return list, total, err
}

func (r *ProductAdminRepository) GetDetail(ctx context.Context, id, shopID uint64) (*model.Product, []model.ProductSku, []model.ProductImage, []model.ProductAttr, error) {
	var p model.Product
	if err := r.conn.QueryRowCtx(ctx, &p, "SELECT "+productColumns+" FROM products WHERE id=? AND shop_id=? LIMIT 1", id, shopID); err != nil {
		return nil, nil, nil, nil, err
	}
	var skus []model.ProductSku
	_ = r.conn.QueryRowsCtx(ctx, &skus, "SELECT "+productSkuColumns+" FROM product_skus WHERE product_id=? AND deleted_at IS NULL", id)
	var imgs []model.ProductImage
	_ = r.conn.QueryRowsCtx(ctx, &imgs, "SELECT "+productImageColumns+" FROM product_images WHERE product_id=? ORDER BY sort ASC, id ASC", id)
	var attrs []model.ProductAttr
	_ = r.conn.QueryRowsCtx(ctx, &attrs, "SELECT "+productAttrColumns+" FROM product_attrs WHERE product_id=?", id)
	return &p, skus, imgs, attrs, nil
}

func (r *ProductAdminRepository) SaveProduct(ctx context.Context, shopID, operatorID uint64, id uint64, req types.MerchantProductSaveReq) (*model.Product, error) {
	status := req.Status
	if status == "" {
		status = model.ProductDraft
	}
	if status != model.ProductDraft && status != model.ProductOnSale {
		return nil, errors.New("创建/保存仅允许 draft 或 on_sale")
	}
	pt := req.ProductType
	if pt == "" {
		pt = model.ProductTypePhysical
	}
	pet := req.PetType
	if pet == "" {
		pet = "both"
	}

	specBytes, _ := json.Marshal(req.SpecJSON)
	var product *model.Product

	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if id == 0 {
			product = &model.Product{
				ShopID:           shopID,
				ProductNo:        fmt.Sprintf("P%d", time.Now().UnixNano()%1e12),
				Name:             req.Name,
				Subtitle:         req.Subtitle,
				Description:      req.Description,
				CategoryID:       req.CategoryID,
				ProductType:      pt,
				SpecJSON:         string(specBytes),
				MainImage:        req.MainImage,
				Status:           status,
				PetType:          pet,
				Discount:         100,
				ShelfLife:        req.ShelfLife,
				StorageCondition: req.StorageCondition,
				StockWarn:        10,
			}
			if status == model.ProductOnSale {
				product.PublishTime = common.LocalTime(time.Now())
			}
			newID, err := lastInsertID(ctx, session,
				`INSERT INTO products (shop_id, product_no, name, subtitle, description, category_id, product_type, spec_json, main_image, status, pet_type, discount, shelf_life, storage_condition, stock_warn, publish_time)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				product.ShopID, product.ProductNo, product.Name, product.Subtitle, product.Description,
				product.CategoryID, product.ProductType, product.SpecJSON, product.MainImage, product.Status,
				product.PetType, product.Discount, product.ShelfLife, product.StorageCondition, product.StockWarn, product.PublishTime,
			)
			if err != nil {
				return err
			}
			product.ID = newID
		} else {
			var existing model.Product
			if err := session.QueryRowCtx(ctx, &existing, "SELECT "+productColumns+" FROM products WHERE id=? AND shop_id=? LIMIT 1", id, shopID); err != nil {
				return errors.New("商品不存在")
			}
			if existing.Status == model.ProductDeleted {
				return errors.New("回收站商品不可编辑")
			}
			product = &existing
			updates := map[string]interface{}{
				"name":              req.Name,
				"subtitle":          req.Subtitle,
				"description":       req.Description,
				"category_id":       req.CategoryID,
				"product_type":      pt,
				"spec_json":         string(specBytes),
				"main_image":        req.MainImage,
				"pet_type":          pet,
				"shelf_life":        req.ShelfLife,
				"storage_condition": req.StorageCondition,
			}
			if status == model.ProductOnSale || status == model.ProductDraft {
				updates["status"] = status
				if status == model.ProductOnSale {
					updates["publish_time"] = time.Now()
				}
			}
			query, args, err := buildUpdate("products", updates, "id=?", id)
			if err != nil {
				return err
			}
			if _, err := session.ExecCtx(ctx, query, args...); err != nil {
				return err
			}
			product.Name = req.Name
			product.Status = status
		}

		if err := r.syncSKUs(ctx, session, product, req); err != nil {
			return err
		}
		if err := r.syncImages(ctx, session, product.ID, shopID, req.Images, req.MainImage); err != nil {
			return err
		}
		if err := r.syncAttrs(ctx, session, product.ID, req.Attrs); err != nil {
			return err
		}
		if err := r.syncTags(ctx, session, product.ID, req.TagIDs); err != nil {
			return err
		}
		return r.aggregateFromSKUs(ctx, session, product.ID)
	})
	if err != nil {
		return nil, err
	}
	action := "create"
	if id > 0 {
		action = "update"
	}
	after, _ := json.Marshal(map[string]interface{}{
		"name": product.Name, "status": product.Status, "sale_price": product.SalePrice,
		"stock": product.Stock, "category_id": product.CategoryID, "product_type": product.ProductType,
	})
	_ = r.AddOpLog(ctx, shopID, &product.ID, operatorID, action, "", string(after))
	return product, nil
}

func (r *ProductAdminRepository) syncSKUs(ctx context.Context, session sqlx.Session, product *model.Product, req types.MerchantProductSaveReq) error {
	skuInputs := req.Skus
	if len(skuInputs) == 0 {
		specs := make([]skuutil.SpecItem, 0, len(req.SpecJSON))
		for _, s := range req.SpecJSON {
			specs = append(specs, skuutil.SpecItem{Name: s.Name, Values: s.Values})
		}
		combos := skuutil.Cartesian(specs)
		price := req.SalePrice
		if price <= 0 {
			price = 0.01
		}
		stock := req.Stock
		for _, c := range combos {
			skuInputs = append(skuInputs, types.SkuInput{
				SpecValues: c,
				SalePrice:  price,
				Stock:      stock,
				StockWarn:  10,
				Status:     model.SKUEnabled,
			})
		}
	}

	var existing []model.ProductSku
	_ = session.QueryRowsCtx(ctx, &existing, "SELECT "+productSkuColumns+" FROM product_skus WHERE product_id=? AND deleted_at IS NULL", product.ID)
	byKey := map[string]model.ProductSku{}
	for _, e := range existing {
		byKey[e.SpecKey] = e
	}
	keepKeys := map[string]struct{}{}

	for _, in := range skuInputs {
		key := skuutil.SpecKey(in.SpecValues)
		keepKeys[key] = struct{}{}
		sv, _ := skuutil.SpecValuesJSON(in.SpecValues)
		st := in.Status
		if st == "" {
			st = model.SKUEnabled
		}
		warn := in.StockWarn
		if warn <= 0 {
			warn = 10
		}
		if old, ok := byKey[key]; ok {
			query, args, err := buildUpdate("product_skus", map[string]interface{}{
				"spec_values": sv, "sale_price": in.SalePrice, "market_price": in.MarketPrice,
				"cost_price": in.CostPrice, "stock": in.Stock, "stock_warn": warn,
				"barcode": in.Barcode, "status": st,
			}, "id=?", old.ID)
			if err != nil {
				return err
			}
			if _, err := session.ExecCtx(ctx, query, args...); err != nil {
				return err
			}
			continue
		}
		_, err := lastInsertID(ctx, session,
			`INSERT INTO product_skus (product_id, shop_id, sku_no, spec_values, spec_key, sale_price, market_price, cost_price, stock, stock_warn, barcode, status)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			product.ID, product.ShopID, skuutil.SkuNo(product.ProductNo, key), sv, key,
			in.SalePrice, in.MarketPrice, in.CostPrice, in.Stock, warn, in.Barcode, st,
		)
		if err != nil {
			return err
		}
	}

	for key, old := range byKey {
		if _, ok := keepKeys[key]; ok {
			continue
		}
		if old.SoldCount > 0 {
			return fmt.Errorf("SKU「%s」已有销量，不可删除，请禁用", old.SkuNo)
		}
		now := common.LocalTime(time.Now())
		query, args, err := buildUpdate("product_skus", map[string]interface{}{
			"deleted_at": now, "status": model.SKUDisabled,
		}, "id=?", old.ID)
		if err != nil {
			return err
		}
		if _, err := session.ExecCtx(ctx, query, args...); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductAdminRepository) syncImages(ctx context.Context, session sqlx.Session, productID, shopID uint64, images []types.ImageInput, mainImage string) error {
	if _, err := session.ExecCtx(ctx, "DELETE FROM product_images WHERE product_id=?", productID); err != nil {
		return err
	}
	if mainImage != "" {
		images = append([]types.ImageInput{{URL: mainImage, Typ: "main", Sort: 0}}, images...)
	}
	for i, im := range images {
		typ := im.Typ
		if typ == "" {
			typ = "gallery"
		}
		sort := im.Sort
		if sort == 0 {
			sort = i
		}
		if _, err := lastInsertID(ctx, session,
			"INSERT INTO product_images (product_id, shop_id, url, typ, sort) VALUES (?, ?, ?, ?, ?)",
			productID, shopID, im.URL, typ, sort,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductAdminRepository) syncAttrs(ctx context.Context, session sqlx.Session, productID uint64, attrs []types.AttrInput) error {
	_, _ = session.ExecCtx(ctx, "DELETE FROM product_attrs WHERE product_id=?", productID)
	for _, a := range attrs {
		if a.AttrKey == "" {
			continue
		}
		if _, err := lastInsertID(ctx, session,
			"INSERT INTO product_attrs (product_id, template_id, attr_key, attr_label, attr_value) VALUES (?, ?, ?, ?, ?)",
			productID, a.TemplateID, a.AttrKey, a.AttrLabel, a.AttrValue,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductAdminRepository) syncTags(ctx context.Context, session sqlx.Session, productID uint64, tagIDs []uint64) error {
	_, _ = session.ExecCtx(ctx, "DELETE FROM product_tag_rels WHERE product_id=?", productID)
	for _, tid := range tagIDs {
		if _, err := session.ExecCtx(ctx, "INSERT INTO product_tag_rels (product_id, tag_id) VALUES (?, ?)", productID, tid); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductAdminRepository) aggregateFromSKUs(ctx context.Context, session sqlx.Session, productID uint64) error {
	var skus []model.ProductSku
	if err := session.QueryRowsCtx(ctx, &skus,
		"SELECT "+productSkuColumns+" FROM product_skus WHERE product_id=? AND deleted_at IS NULL AND status=?",
		productID, model.SKUEnabled,
	); err != nil {
		return err
	}
	totalStock := 0
	minPrice := 0.0
	for i, s := range skus {
		totalStock += s.Stock
		if i == 0 || s.SalePrice < minPrice {
			minPrice = s.SalePrice
		}
	}
	query, args, err := buildUpdate("products", map[string]interface{}{
		"stock": totalStock, "sale_price": minPrice,
	}, "id=?", productID)
	if err != nil {
		return err
	}
	_, err = session.ExecCtx(ctx, query, args...)
	return err
}

func (r *ProductAdminRepository) SetStatus(ctx context.Context, id, shopID uint64, status string) error {
	updates := map[string]interface{}{"status": status}
	if status == model.ProductOnSale {
		updates["publish_time"] = time.Now()
	}
	if status == model.ProductDeleted {
		updates["deleted_at"] = time.Now()
	}
	if status != model.ProductDeleted {
		updates["deleted_at"] = nil
	}
	query, args, err := buildUpdate("products", updates, "id=? AND shop_id=?", id, shopID)
	if err != nil {
		return err
	}
	n, err := execAffected(ctx, r.conn, query, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("商品不存在")
	}
	return nil
}

func (r *ProductAdminRepository) GetByID(ctx context.Context, id uint64) (*model.Product, error) {
	var p model.Product
	if err := r.conn.QueryRowCtx(ctx, &p, "SELECT "+productColumns+" FROM products WHERE id=? LIMIT 1", id); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductAdminRepository) SetStatusByID(ctx context.Context, id uint64, status string) error {
	p, err := r.GetByID(ctx, id)
	if err != nil {
		return errors.New("商品不存在")
	}
	return r.SetStatus(ctx, id, p.ShopID, status)
}

func (r *ProductAdminRepository) CopyProduct(ctx context.Context, id, shopID, operatorID uint64) (*model.Product, error) {
	p, skus, imgs, attrs, err := r.GetDetail(ctx, id, shopID)
	if err != nil {
		return nil, err
	}
	var specs []types.SpecItem
	_ = json.Unmarshal([]byte(p.SpecJSON), &specs)
	skuIns := make([]types.SkuInput, 0, len(skus))
	for _, s := range skus {
		var m map[string]string
		_ = json.Unmarshal([]byte(s.SpecValues), &m)
		skuIns = append(skuIns, types.SkuInput{
			SpecValues: m, SalePrice: s.SalePrice, MarketPrice: s.MarketPrice,
			CostPrice: s.CostPrice, Stock: s.Stock, StockWarn: s.StockWarn, Status: model.SKUEnabled,
		})
	}
	imgIns := make([]types.ImageInput, 0, len(imgs))
	for _, im := range imgs {
		if im.Typ == "main" {
			continue
		}
		imgIns = append(imgIns, types.ImageInput{URL: im.URL, Typ: im.Typ, Sort: im.Sort})
	}
	attrIns := make([]types.AttrInput, 0, len(attrs))
	for _, a := range attrs {
		attrIns = append(attrIns, types.AttrInput{TemplateID: a.TemplateID, AttrKey: a.AttrKey, AttrLabel: a.AttrLabel, AttrValue: a.AttrValue})
	}
	req := types.MerchantProductSaveReq{
		Name: p.Name + "（副本）", Subtitle: p.Subtitle, Description: p.Description,
		CategoryID: p.CategoryID, ProductType: p.ProductType, PetType: p.PetType,
		Status: model.ProductDraft, MainImage: p.MainImage, SpecJSON: specs, Skus: skuIns,
		Images: imgIns, Attrs: attrIns, ShelfLife: p.ShelfLife, StorageCondition: p.StorageCondition,
		SalePrice: p.SalePrice, Stock: p.Stock,
	}
	np, err := r.SaveProduct(ctx, shopID, operatorID, 0, req)
	if err != nil {
		return nil, err
	}
	_, _ = r.conn.ExecCtx(ctx, "UPDATE products SET copy_from_id=? WHERE id=?", id, np.ID)
	return np, nil
}

func (r *ProductAdminRepository) PermanentDelete(ctx context.Context, shopID uint64, ids []uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, id := range ids {
			var p model.Product
			if err := session.QueryRowCtx(ctx, &p, "SELECT "+productColumns+" FROM products WHERE id=? AND shop_id=? AND status=? LIMIT 1", id, shopID, model.ProductDeleted); err != nil {
				return fmt.Errorf("商品 %d 不在回收站", id)
			}
			_, _ = session.ExecCtx(ctx, "DELETE FROM product_skus WHERE product_id=?", id)
			_, _ = session.ExecCtx(ctx, "DELETE FROM product_images WHERE product_id=?", id)
			_, _ = session.ExecCtx(ctx, "DELETE FROM product_attrs WHERE product_id=?", id)
			_, _ = session.ExecCtx(ctx, "DELETE FROM product_tag_rels WHERE product_id=?", id)
			_, _ = session.ExecCtx(ctx, "DELETE FROM product_schedules WHERE product_id=?", id)
			if _, err := session.ExecCtx(ctx, "DELETE FROM products WHERE id=?", id); err != nil {
				return err
			}
			_ = p
		}
		return nil
	})
}

func (r *ProductAdminRepository) AdjustSkuStock(ctx context.Context, shopID uint64, req types.StockAdjustReq) error {
	var sku model.ProductSku
	if err := r.conn.QueryRowCtx(ctx, &sku, "SELECT "+productSkuColumns+" FROM product_skus WHERE id=? AND shop_id=? AND deleted_at IS NULL LIMIT 1", req.SkuID, shopID); err != nil {
		return errors.New("SKU不存在")
	}
	newStock := sku.Stock
	if req.Stock != nil {
		newStock = *req.Stock
	} else if req.Delta != nil {
		newStock = sku.Stock + *req.Delta
	}
	if newStock < 0 {
		return errors.New("库存不能为负")
	}
	if _, err := r.conn.ExecCtx(ctx, "UPDATE product_skus SET stock=? WHERE id=?", newStock, sku.ID); err != nil {
		return err
	}
	return r.aggregateFromSKUs(ctx, r.conn, sku.ProductID)
}

func (r *ProductAdminRepository) ListStockWarnings(ctx context.Context, shopID uint64, page, pageSize int) ([]model.ProductSku, int64, error) {
	where := "shop_id=? AND deleted_at IS NULL AND stock <= stock_warn"
	total, _ := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM product_skus WHERE "+where, shopID)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var list []model.ProductSku
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+productSkuColumns+" FROM product_skus WHERE "+where+" ORDER BY stock ASC LIMIT ? OFFSET ?",
		shopID, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *ProductAdminRepository) CreateBatchJob(ctx context.Context, job *model.ProductBatchJob) error {
	id, err := lastInsertID(ctx, r.conn,
		`INSERT INTO product_batch_jobs (shop_id, job_type, payload_json, progress, total, status, result_msg, operator_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		job.ShopID, job.JobType, job.PayloadJSON, job.Progress, job.Total, job.Status, job.ResultMsg, job.OperatorID,
	)
	if err != nil {
		return err
	}
	job.ID = id
	return nil
}

func (r *ProductAdminRepository) GetBatchJob(ctx context.Context, id, shopID uint64) (*model.ProductBatchJob, error) {
	var j model.ProductBatchJob
	err := r.conn.QueryRowCtx(ctx, &j, "SELECT "+productBatchJobColumns+" FROM product_batch_jobs WHERE id=? AND shop_id=? LIMIT 1", id, shopID)
	return &j, err
}

func (r *ProductAdminRepository) UpdateBatchJob(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("product_batch_jobs", updates, "id=?", id)
	if err != nil {
		return err
	}
	_, err = r.conn.ExecCtx(ctx, query, args...)
	return err
}

func (r *ProductAdminRepository) UpdateProductCategory(ctx context.Context, id, shopID, categoryID uint64) error {
	_, err := r.conn.ExecCtx(ctx, "UPDATE products SET category_id=? WHERE id=? AND shop_id=?", categoryID, id, shopID)
	return err
}

func (r *ProductAdminRepository) ListSkusByProductShop(ctx context.Context, productID, shopID uint64) ([]model.ProductSku, error) {
	var skus []model.ProductSku
	err := r.conn.QueryRowsCtx(ctx, &skus,
		"SELECT "+productSkuColumns+" FROM product_skus WHERE product_id=? AND shop_id=? AND deleted_at IS NULL",
		productID, shopID,
	)
	return skus, err
}

func (r *ProductAdminRepository) UpdateSkuSalePrice(ctx context.Context, skuID uint64, price float64) error {
	_, err := r.conn.ExecCtx(ctx, "UPDATE product_skus SET sale_price=? WHERE id=?", price, skuID)
	return err
}

func (r *ProductAdminRepository) AddOpLog(ctx context.Context, shopID uint64, productID *uint64, operatorID uint64, action, before, after string) error {
	before = jsonOrNull(before)
	after = jsonOrNull(after)
	_, err := lastInsertID(ctx, r.conn,
		`INSERT INTO product_op_logs (shop_id, product_id, operator_id, action, before_json, after_json)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		shopID, productID, operatorID, action, before, after,
	)
	return err
}

func jsonOrNull(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "null"
	}
	if json.Valid([]byte(s)) {
		return s
	}
	b, _ := json.Marshal(s)
	return string(b)
}

type OpLogItem struct {
	ID           uint64           `json:"id"`
	ShopID       uint64           `json:"shop_id"`
	ProductID    *uint64          `json:"product_id,omitempty"`
	ProductName  string           `json:"product_name"`
	OperatorID   uint64           `json:"operator_id"`
	OperatorName string           `json:"operator_name"`
	Action       string           `json:"action"`
	ActionLabel  string           `json:"action_label"`
	TargetType   string           `json:"target_type"`
	TargetName   string           `json:"target_name"`
	Summary      string           `json:"summary"`
	BeforeJSON   string           `json:"before_json,omitempty"`
	AfterJSON    string           `json:"after_json,omitempty"`
	CreatedAt    common.LocalTime `json:"created_at"`
}

func (r *ProductAdminRepository) ListOpLogs(ctx context.Context, shopID uint64, productID uint64, page, pageSize int) ([]OpLogItem, int64, error) {
	where := []string{"shop_id=?"}
	args := []any{shopID}
	if productID > 0 {
		where = append(where, "product_id=?")
		args = append(args, productID)
	}
	whereSQL := strings.Join(where, " AND ")

	total, _ := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM product_op_logs WHERE "+whereSQL, args...)
	var logs []model.ProductOpLog
	qArgs := append(args, pageSize, (page-1)*pageSize)
	err := r.conn.QueryRowsCtx(ctx, &logs,
		"SELECT "+productOpLogColumns+" FROM product_op_logs WHERE "+whereSQL+" ORDER BY id DESC LIMIT ? OFFSET ?",
		qArgs...,
	)
	if err != nil {
		return nil, 0, err
	}

	productIDs := make([]uint64, 0)
	operatorIDs := make([]uint64, 0)
	seenP, seenO := map[uint64]struct{}{}, map[uint64]struct{}{}
	for _, lg := range logs {
		if lg.ProductID != nil && *lg.ProductID > 0 {
			if _, ok := seenP[*lg.ProductID]; !ok {
				seenP[*lg.ProductID] = struct{}{}
				productIDs = append(productIDs, *lg.ProductID)
			}
		}
		if lg.OperatorID > 0 {
			if _, ok := seenO[lg.OperatorID]; !ok {
				seenO[lg.OperatorID] = struct{}{}
				operatorIDs = append(operatorIDs, lg.OperatorID)
			}
		}
	}

	nameByProduct := map[uint64]string{}
	if len(productIDs) > 0 {
		type row struct {
			ID   uint64 `db:"id"`
			Name string `db:"name"`
		}
		var rows []row
		_ = r.conn.QueryRowsCtx(ctx, &rows, "SELECT id, name FROM products WHERE id IN ("+placeholders(len(productIDs))+")", inArgs(productIDs)...)
		for _, x := range rows {
			nameByProduct[x.ID] = x.Name
		}
	}
	nameByUser := map[uint64]string{}
	if len(operatorIDs) > 0 {
		type row struct {
			ID       uint64 `db:"id"`
			Nickname string `db:"nickname"`
			Mobile   string `db:"mobile"`
		}
		var rows []row
		_ = r.conn.QueryRowsCtx(ctx, &rows, "SELECT id, nickname, mobile FROM users WHERE id IN ("+placeholders(len(operatorIDs))+")", inArgs(operatorIDs)...)
		for _, x := range rows {
			n := strings.TrimSpace(x.Nickname)
			if n == "" {
				n = x.Mobile
			}
			nameByUser[x.ID] = n
		}
	}

	out := make([]OpLogItem, 0, len(logs))
	for _, lg := range logs {
		item := OpLogItem{
			ID: lg.ID, ShopID: lg.ShopID, ProductID: lg.ProductID,
			OperatorID: lg.OperatorID, Action: lg.Action,
			BeforeJSON: lg.BeforeJSON, AfterJSON: lg.AfterJSON, CreatedAt: lg.CreatedAt,
			TargetType: "商品", ActionLabel: opActionLabel(lg.Action),
		}
		if lg.OperatorID > 0 {
			item.OperatorName = nameByUser[lg.OperatorID]
		}
		if item.OperatorName == "" && lg.OperatorID > 0 {
			item.OperatorName = fmt.Sprintf("用户#%d", lg.OperatorID)
		}
		if lg.ProductID != nil && *lg.ProductID > 0 {
			item.ProductName = nameByProduct[*lg.ProductID]
			if item.ProductName == "" {
				item.ProductName = pickJSONString(lg.AfterJSON, "name")
			}
			if item.ProductName == "" {
				item.ProductName = fmt.Sprintf("商品#%d", *lg.ProductID)
			}
			item.TargetName = item.ProductName
		} else if lg.Action == "permanent_delete" {
			item.TargetName = "批量商品"
			item.ProductName = summarizeDeletedIDs(lg.AfterJSON)
		} else {
			item.TargetName = "-"
		}
		item.Summary = buildOpSummary(item)
		out = append(out, item)
	}
	return out, total, nil
}

func opActionLabel(action string) string {
	m := map[string]string{
		"create": "创建商品", "update": "更新商品", "copy": "复制商品",
		"schedule": "设置定时", "permanent_delete": "永久删除",
		"status:on_sale": "上架", "status:off_sale": "下架",
		"status:deleted": "移入回收站", "status:draft": "设为草稿",
		"save":              "保存商品",
		"platform_off_sale": "平台强制下架",
		"platform_delete":   "平台删除",
	}
	if v, ok := m[action]; ok {
		return v
	}
	if strings.HasPrefix(action, "status:") {
		return "变更状态"
	}
	return action
}

func pickJSONString(raw, key string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return ""
	}
	if v, ok := m[key]; ok {
		switch t := v.(type) {
		case string:
			return t
		default:
			return fmt.Sprint(t)
		}
	}
	return ""
}

func summarizeDeletedIDs(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return "已删除商品"
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return "已删除商品"
	}
	ids, _ := m["deleted_ids"].([]interface{})
	if len(ids) == 0 {
		return "已删除商品"
	}
	return fmt.Sprintf("共 %d 件商品", len(ids))
}

func buildOpSummary(item OpLogItem) string {
	name := item.TargetName
	if name == "" || name == "-" {
		name = "商品"
	}
	switch {
	case item.Action == "create":
		return fmt.Sprintf("创建了「%s」", name)
	case item.Action == "update":
		return fmt.Sprintf("更新了「%s」", name)
	case item.Action == "copy":
		return fmt.Sprintf("复制生成了「%s」", name)
	case item.Action == "status:on_sale":
		return fmt.Sprintf("将「%s」上架", name)
	case item.Action == "status:off_sale":
		return fmt.Sprintf("将「%s」下架", name)
	case item.Action == "status:deleted":
		return fmt.Sprintf("将「%s」移入回收站", name)
	case item.Action == "status:draft":
		return fmt.Sprintf("将「%s」设为草稿", name)
	case item.Action == "schedule":
		act := pickJSONString(item.AfterJSON, "action")
		runAt := pickJSONString(item.AfterJSON, "run_at")
		actCN := map[string]string{"on_sale": "上架", "off_sale": "下架"}[act]
		if actCN == "" {
			actCN = "变更状态"
		}
		if runAt != "" {
			return fmt.Sprintf("为「%s」设置定时%s（%s）", name, actCN, runAt)
		}
		return fmt.Sprintf("为「%s」设置定时%s", name, actCN)
	case item.Action == "permanent_delete":
		return fmt.Sprintf("永久删除了%s", name)
	case item.Action == "platform_off_sale":
		if remark := pickJSONString(item.AfterJSON, "remark"); remark != "" {
			return fmt.Sprintf("平台强制下架「%s」，备注：%s", name, remark)
		}
		return fmt.Sprintf("平台强制下架「%s」", name)
	case item.Action == "platform_delete":
		if remark := pickJSONString(item.AfterJSON, "remark"); remark != "" {
			return fmt.Sprintf("平台将「%s」移入回收站，备注：%s", name, remark)
		}
		return fmt.Sprintf("平台将「%s」移入回收站", name)
	default:
		return fmt.Sprintf("%s：%s", item.ActionLabel, name)
	}
}

func (r *ProductAdminRepository) ReserveSkuStock(ctx context.Context, items []SkuStockItem) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, it := range items {
			if it.SkuID > 0 {
				n, err := execSession(ctx, session,
					"UPDATE product_skus SET stock=stock-? WHERE id=? AND deleted_at IS NULL AND stock>=?",
					it.Quantity, it.SkuID, it.Quantity,
				)
				if err != nil {
					return err
				}
				if n == 0 {
					return fmt.Errorf("SKU %d 库存不足", it.SkuID)
				}
				var sku model.ProductSku
				_ = session.QueryRowCtx(ctx, &sku, "SELECT product_id FROM product_skus WHERE id=? LIMIT 1", it.SkuID)
				_ = r.aggregateFromSKUs(ctx, session, sku.ProductID)
				continue
			}
			skuID := r.firstSkuIDSession(ctx, session, it.ProductID)
			if skuID > 0 {
				n, err := execSession(ctx, session,
					"UPDATE product_skus SET stock=stock-? WHERE id=? AND deleted_at IS NULL AND stock>=?",
					it.Quantity, skuID, it.Quantity,
				)
				if err != nil {
					return err
				}
				if n == 0 {
					return fmt.Errorf("SKU %d 库存不足", skuID)
				}
				_ = r.aggregateFromSKUs(ctx, session, it.ProductID)
				continue
			}
			n, err := execSession(ctx, session,
				"UPDATE products SET stock=stock-? WHERE id=? AND stock>=?", it.Quantity, it.ProductID, it.Quantity)
			if err != nil {
				return err
			}
			if n == 0 {
				return fmt.Errorf("商品 %d 库存不足", it.ProductID)
			}
		}
		return nil
	})
}

func (r *ProductAdminRepository) ReleaseSkuStock(ctx context.Context, items []SkuStockItem) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, it := range items {
			skuID := it.SkuID
			if skuID == 0 {
				skuID = r.firstSkuIDSession(ctx, session, it.ProductID)
			}
			if skuID > 0 {
				if _, err := session.ExecCtx(ctx, "UPDATE product_skus SET stock=stock+? WHERE id=?", it.Quantity, skuID); err != nil {
					return err
				}
				var sku model.ProductSku
				_ = session.QueryRowCtx(ctx, &sku, "SELECT product_id FROM product_skus WHERE id=? LIMIT 1", skuID)
				_ = r.aggregateFromSKUs(ctx, session, sku.ProductID)
				continue
			}
			_, _ = session.ExecCtx(ctx, "UPDATE products SET stock=stock+? WHERE id=?", it.Quantity, it.ProductID)
		}
		return nil
	})
}

type SkuStockItem struct {
	ProductID uint64
	SkuID     uint64
	Quantity  int
}

func (r *ProductAdminRepository) ListTags(ctx context.Context, shopID uint64) ([]model.ProductTag, error) {
	var list []model.ProductTag
	err := r.conn.QueryRowsCtx(ctx, &list, "SELECT "+productTagColumns+" FROM product_tags WHERE shop_id IN (0, ?) ORDER BY id DESC", shopID)
	return list, err
}

func (r *ProductAdminRepository) SaveTag(ctx context.Context, tag *model.ProductTag) error {
	if tag.ID == 0 {
		id, err := lastInsertID(ctx, r.conn, "INSERT INTO product_tags (shop_id, name, color, status) VALUES (?, ?, ?, ?)",
			tag.ShopID, tag.Name, tag.Color, tag.Status)
		if err != nil {
			return err
		}
		tag.ID = id
		return nil
	}
	query, args, err := buildUpdate("product_tags", map[string]interface{}{"name": tag.Name, "color": tag.Color, "status": tag.Status}, "id=?", tag.ID)
	if err != nil {
		return err
	}
	_, err = r.conn.ExecCtx(ctx, query, args...)
	return err
}

func (r *ProductAdminRepository) DeleteTag(ctx context.Context, id, shopID uint64) error {
	_, err := r.conn.ExecCtx(ctx, "DELETE FROM product_tags WHERE id=? AND shop_id=?", id, shopID)
	return err
}

func (r *ProductAdminRepository) ListAttrTemplates(ctx context.Context, shopID uint64) ([]model.ProductAttrTemplate, error) {
	var list []model.ProductAttrTemplate
	err := r.conn.QueryRowsCtx(ctx, &list, "SELECT "+productAttrTemplateColumns+" FROM product_attr_templates WHERE shop_id IN (0, ?) ORDER BY id DESC", shopID)
	return list, err
}

func (r *ProductAdminRepository) SaveAttrTemplate(ctx context.Context, t *model.ProductAttrTemplate) error {
	if t.ID == 0 {
		id, err := lastInsertID(ctx, r.conn, "INSERT INTO product_attr_templates (shop_id, name, attrs_json, status) VALUES (?, ?, ?, ?)",
			t.ShopID, t.Name, t.AttrsJSON, t.Status)
		if err != nil {
			return err
		}
		t.ID = id
		return nil
	}
	query, args, err := buildUpdate("product_attr_templates", map[string]interface{}{"name": t.Name, "attrs_json": t.AttrsJSON, "status": t.Status}, "id=?", t.ID)
	if err != nil {
		return err
	}
	_, err = r.conn.ExecCtx(ctx, query, args...)
	return err
}

func (r *ProductAdminRepository) DeleteAttrTemplate(ctx context.Context, id, shopID uint64) error {
	_, err := r.conn.ExecCtx(ctx, "DELETE FROM product_attr_templates WHERE id=? AND shop_id=?", id, shopID)
	return err
}

func (r *ProductAdminRepository) CreateSchedule(ctx context.Context, s *model.ProductSchedule) error {
	_, _ = r.conn.ExecCtx(ctx,
		"UPDATE product_schedules SET status='cancelled' WHERE product_id=? AND action=? AND status='pending'",
		s.ProductID, s.Action)
	id, err := lastInsertID(ctx, r.conn,
		"INSERT INTO product_schedules (product_id, shop_id, action, run_at, status) VALUES (?, ?, ?, ?, ?)",
		s.ProductID, s.ShopID, s.Action, s.RunAt, s.Status,
	)
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (r *ProductAdminRepository) CancelSchedule(ctx context.Context, id, shopID uint64) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE product_schedules SET status='cancelled' WHERE id=? AND shop_id=? AND status='pending'", id, shopID)
	return err
}

func (r *ProductAdminRepository) ClaimDueSchedules(ctx context.Context, limit int) ([]model.ProductSchedule, error) {
	var list []model.ProductSchedule
	now := time.Now()
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+productScheduleColumns+" FROM product_schedules WHERE status='pending' AND run_at<=? AND locked_at IS NULL LIMIT ?",
		now, limit,
	)
	if err != nil {
		return nil, err
	}
	out := make([]model.ProductSchedule, 0, len(list))
	for _, s := range list {
		n, err := execAffected(ctx, r.conn,
			"UPDATE product_schedules SET locked_at=? WHERE id=? AND status='pending' AND locked_at IS NULL", now, s.ID)
		if err != nil {
			continue
		}
		if n == 1 {
			out = append(out, s)
		}
	}
	return out, nil
}

func (r *ProductAdminRepository) FinishSchedule(ctx context.Context, id uint64, ok bool) error {
	st := "done"
	if !ok {
		st = "cancelled"
	}
	_, err := r.conn.ExecCtx(ctx, "UPDATE product_schedules SET status=? WHERE id=?", st, id)
	return err
}

func (r *ProductAdminRepository) FirstSkuID(ctx context.Context, productID uint64) uint64 {
	return r.firstSkuIDSession(ctx, r.conn, productID)
}

func (r *ProductAdminRepository) firstSkuIDSession(ctx context.Context, session sqlx.Session, productID uint64) uint64 {
	var id uint64
	_ = session.QueryRowCtx(ctx, &id,
		"SELECT id FROM product_skus WHERE product_id=? AND deleted_at IS NULL ORDER BY id ASC LIMIT 1", productID)
	return id
}

func (r *ProductAdminRepository) GetSku(ctx context.Context, id uint64) (*model.ProductSku, error) {
	var s model.ProductSku
	err := r.conn.QueryRowCtx(ctx, &s, "SELECT "+productSkuColumns+" FROM product_skus WHERE id=? AND deleted_at IS NULL LIMIT 1", id)
	return &s, err
}

func (r *ProductAdminRepository) AggregatePublic(ctx context.Context, productID uint64) error {
	return r.aggregateFromSKUs(ctx, r.conn, productID)
}
