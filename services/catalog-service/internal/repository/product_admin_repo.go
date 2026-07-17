package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/model"
	"mymall/services/catalog-service/internal/skuutil"
	"mymall/services/catalog-service/internal/types"

	"gorm.io/gorm"
)

type ProductAdminRepository struct {
	db *gorm.DB
}

func NewProductAdminRepository(db *gorm.DB) *ProductAdminRepository {
	return &ProductAdminRepository{db: db}
}

type ProductListFilter struct {
	ShopID    uint64
	Name      string
	ProductNo string
	CategoryID uint64
	Status    string
	ProductType string
	StockWarnOnly bool
	Page      int
	PageSize  int
	OrderBy   string
	Recycle   bool // true=仅回收站
}

func (r *ProductAdminRepository) List(f ProductListFilter) ([]model.Product, int64, error) {
	q := r.db.Model(&model.Product{}).Where("shop_id = ?", f.ShopID)
	if f.Recycle {
		q = q.Where("status = ?", model.ProductDeleted)
	} else {
		q = q.Where("status <> ?", model.ProductDeleted)
	}
	if f.Name != "" {
		q = q.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.ProductNo != "" {
		q = q.Where("product_no LIKE ?", "%"+f.ProductNo+"%")
	}
	if f.CategoryID > 0 {
		q = q.Where("category_id = ?", f.CategoryID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.ProductType != "" {
		q = q.Where("product_type = ?", f.ProductType)
	}
	if f.StockWarnOnly {
		q = q.Where("stock <= stock_warn")
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
	case "created_asc":
		order = "id ASC"
	}
	var list []model.Product
	err := q.Order(order).Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&list).Error
	return list, total, err
}

func (r *ProductAdminRepository) GetDetail(id, shopID uint64) (*model.Product, []model.ProductSku, []model.ProductImage, []model.ProductAttr, error) {
	var p model.Product
	if err := r.db.Where("id = ? AND shop_id = ?", id, shopID).First(&p).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	var skus []model.ProductSku
	_ = r.db.Where("product_id = ? AND deleted_at IS NULL", id).Find(&skus).Error
	var imgs []model.ProductImage
	_ = r.db.Where("product_id = ?", id).Order("sort ASC, id ASC").Find(&imgs).Error
	var attrs []model.ProductAttr
	_ = r.db.Where("product_id = ?", id).Find(&attrs).Error
	return &p, skus, imgs, attrs, nil
}

// SaveProduct 创建或更新商品 + SKU/图片/参数（事务）
func (r *ProductAdminRepository) SaveProduct(shopID, operatorID uint64, id uint64, req types.MerchantProductSaveReq) (*model.Product, error) {
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

	err := r.db.Transaction(func(tx *gorm.DB) error {
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
			if err := tx.Create(product).Error; err != nil {
				return err
			}
		} else {
			var existing model.Product
			if err := tx.Where("id = ? AND shop_id = ?", id, shopID).First(&existing).Error; err != nil {
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
			if err := tx.Model(&existing).Updates(updates).Error; err != nil {
				return err
			}
			product.Name = req.Name
			product.Status = status
		}

		if err := r.syncSKUs(tx, product, req); err != nil {
			return err
		}
		if err := r.syncImages(tx, product.ID, shopID, req.Images, req.MainImage); err != nil {
			return err
		}
		if err := r.syncAttrs(tx, product.ID, req.Attrs); err != nil {
			return err
		}
		if err := r.syncTags(tx, product.ID, req.TagIDs); err != nil {
			return err
		}
		// 聚合主表价库存
		return r.aggregateFromSKUs(tx, product.ID)
	})
	if err != nil {
		return nil, err
	}
	_ = r.AddOpLog(shopID, &product.ID, operatorID, "save", "", "")
	return product, nil
}

func (r *ProductAdminRepository) syncSKUs(tx *gorm.DB, product *model.Product, req types.MerchantProductSaveReq) error {
	// 若前端传了 skus 则以其为准；否则按 spec 笛卡尔积生成
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
	_ = tx.Where("product_id = ? AND deleted_at IS NULL", product.ID).Find(&existing).Error
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
			if err := tx.Model(&old).Updates(map[string]interface{}{
				"spec_values":  sv,
				"sale_price":   in.SalePrice,
				"market_price": in.MarketPrice,
				"cost_price":   in.CostPrice,
				"stock":        in.Stock,
				"stock_warn":   warn,
				"barcode":      in.Barcode,
				"status":       st,
			}).Error; err != nil {
				return err
			}
			continue
		}
		sku := model.ProductSku{
			ProductID:   product.ID,
			ShopID:      product.ShopID,
			SkuNo:       skuutil.SkuNo(product.ProductNo, key),
			SpecValues:  sv,
			SpecKey:     key,
			SalePrice:   in.SalePrice,
			MarketPrice: in.MarketPrice,
			CostPrice:   in.CostPrice,
			Stock:       in.Stock,
			StockWarn:   warn,
			Barcode:     in.Barcode,
			Status:      st,
		}
		if err := tx.Create(&sku).Error; err != nil {
			return err
		}
	}

	// 删除未保留的 SKU（已售禁止删）
	for key, old := range byKey {
		if _, ok := keepKeys[key]; ok {
			continue
		}
		if old.SoldCount > 0 {
			return fmt.Errorf("SKU「%s」已有销量，不可删除，请禁用", old.SkuNo)
		}
		now := common.LocalTime(time.Now())
		if err := tx.Model(&old).Updates(map[string]interface{}{
			"deleted_at": now,
			"status":     model.SKUDisabled,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductAdminRepository) syncImages(tx *gorm.DB, productID, shopID uint64, images []types.ImageInput, mainImage string) error {
	if err := tx.Where("product_id = ?", productID).Delete(&model.ProductImage{}).Error; err != nil {
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
		row := model.ProductImage{ProductID: productID, ShopID: shopID, URL: im.URL, Typ: typ, Sort: sort}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductAdminRepository) syncAttrs(tx *gorm.DB, productID uint64, attrs []types.AttrInput) error {
	_ = tx.Where("product_id = ?", productID).Delete(&model.ProductAttr{}).Error
	for _, a := range attrs {
		if a.AttrKey == "" {
			continue
		}
		row := model.ProductAttr{
			ProductID: productID, TemplateID: a.TemplateID,
			AttrKey: a.AttrKey, AttrLabel: a.AttrLabel, AttrValue: a.AttrValue,
		}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductAdminRepository) syncTags(tx *gorm.DB, productID uint64, tagIDs []uint64) error {
	_ = tx.Where("product_id = ?", productID).Delete(&model.ProductTagRel{}).Error
	for _, tid := range tagIDs {
		if err := tx.Create(&model.ProductTagRel{ProductID: productID, TagID: tid}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductAdminRepository) aggregateFromSKUs(tx *gorm.DB, productID uint64) error {
	var skus []model.ProductSku
	if err := tx.Where("product_id = ? AND deleted_at IS NULL AND status = ?", productID, model.SKUEnabled).Find(&skus).Error; err != nil {
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
	return tx.Model(&model.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
		"stock":      totalStock,
		"sale_price": minPrice,
	}).Error
}

func (r *ProductAdminRepository) SetStatus(id, shopID uint64, status string) error {
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
	res := r.db.Model(&model.Product{}).Where("id = ? AND shop_id = ?", id, shopID).Updates(updates)
	if res.RowsAffected == 0 {
		return errors.New("商品不存在")
	}
	return res.Error
}

func (r *ProductAdminRepository) CopyProduct(id, shopID, operatorID uint64) (*model.Product, error) {
	p, skus, imgs, attrs, err := r.GetDetail(id, shopID)
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
	np, err := r.SaveProduct(shopID, operatorID, 0, req)
	if err != nil {
		return nil, err
	}
	_ = r.db.Model(np).Update("copy_from_id", id).Error
	return np, nil
}

func (r *ProductAdminRepository) PermanentDelete(shopID uint64, ids []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var p model.Product
			if err := tx.Where("id = ? AND shop_id = ? AND status = ?", id, shopID, model.ProductDeleted).First(&p).Error; err != nil {
				return fmt.Errorf("商品 %d 不在回收站", id)
			}
			_ = tx.Where("product_id = ?", id).Delete(&model.ProductSku{}).Error
			_ = tx.Where("product_id = ?", id).Delete(&model.ProductImage{}).Error
			_ = tx.Where("product_id = ?", id).Delete(&model.ProductAttr{}).Error
			_ = tx.Where("product_id = ?", id).Delete(&model.ProductTagRel{}).Error
			_ = tx.Where("product_id = ?", id).Delete(&model.ProductSchedule{}).Error
			if err := tx.Delete(&p).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ProductAdminRepository) AdjustSkuStock(shopID uint64, req types.StockAdjustReq) error {
	var sku model.ProductSku
	if err := r.db.Where("id = ? AND shop_id = ? AND deleted_at IS NULL", req.SkuID, shopID).First(&sku).Error; err != nil {
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
	if err := r.db.Model(&sku).Update("stock", newStock).Error; err != nil {
		return err
	}
	return r.aggregateFromSKUs(r.db, sku.ProductID)
}

func (r *ProductAdminRepository) ListStockWarnings(shopID uint64, page, pageSize int) ([]model.ProductSku, int64, error) {
	q := r.db.Model(&model.ProductSku{}).Where("shop_id = ? AND deleted_at IS NULL AND stock <= stock_warn", shopID)
	var total int64
	_ = q.Count(&total)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var list []model.ProductSku
	err := q.Order("stock ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ProductAdminRepository) CreateBatchJob(job *model.ProductBatchJob) error {
	return r.db.Create(job).Error
}

func (r *ProductAdminRepository) GetBatchJob(id, shopID uint64) (*model.ProductBatchJob, error) {
	var j model.ProductBatchJob
	err := r.db.Where("id = ? AND shop_id = ?", id, shopID).First(&j).Error
	return &j, err
}

func (r *ProductAdminRepository) AddOpLog(shopID uint64, productID *uint64, operatorID uint64, action, before, after string) error {
	return r.db.Create(&model.ProductOpLog{
		ShopID: shopID, ProductID: productID, OperatorID: operatorID,
		Action: action, BeforeJSON: before, AfterJSON: after,
	}).Error
}

func (r *ProductAdminRepository) ListOpLogs(shopID uint64, productID uint64, page, pageSize int) ([]model.ProductOpLog, int64, error) {
	q := r.db.Model(&model.ProductOpLog{}).Where("shop_id = ?", shopID)
	if productID > 0 {
		q = q.Where("product_id = ?", productID)
	}
	var total int64
	_ = q.Count(&total)
	var list []model.ProductOpLog
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ProductAdminRepository) ReserveSkuStock(items []SkuStockItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, it := range items {
			if it.SkuID > 0 {
				res := tx.Model(&model.ProductSku{}).
					Where("id = ? AND deleted_at IS NULL AND stock >= ?", it.SkuID, it.Quantity).
					Update("stock", gorm.Expr("stock - ?", it.Quantity))
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					return fmt.Errorf("SKU %d 库存不足", it.SkuID)
				}
				var sku model.ProductSku
				_ = tx.First(&sku, it.SkuID).Error
				_ = r.aggregateFromSKUs(tx, sku.ProductID)
				continue
			}
			// 兼容：无 sku_id 时优先扣默认 SKU，否则扣商品级库存
			skuID := r.FirstSkuID(it.ProductID)
			if skuID > 0 {
				res := tx.Model(&model.ProductSku{}).
					Where("id = ? AND deleted_at IS NULL AND stock >= ?", skuID, it.Quantity).
					Update("stock", gorm.Expr("stock - ?", it.Quantity))
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					return fmt.Errorf("SKU %d 库存不足", skuID)
				}
				_ = r.aggregateFromSKUs(tx, it.ProductID)
				continue
			}
			res := tx.Model(&model.Product{}).
				Where("id = ? AND stock >= ?", it.ProductID, it.Quantity).
				Update("stock", gorm.Expr("stock - ?", it.Quantity))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return fmt.Errorf("商品 %d 库存不足", it.ProductID)
			}
		}
		return nil
	})
}

func (r *ProductAdminRepository) ReleaseSkuStock(items []SkuStockItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, it := range items {
			skuID := it.SkuID
			if skuID == 0 {
				skuID = r.FirstSkuID(it.ProductID)
			}
			if skuID > 0 {
				if err := tx.Model(&model.ProductSku{}).Where("id = ?", skuID).
					Update("stock", gorm.Expr("stock + ?", it.Quantity)).Error; err != nil {
					return err
				}
				var sku model.ProductSku
				_ = tx.First(&sku, skuID).Error
				_ = r.aggregateFromSKUs(tx, sku.ProductID)
				continue
			}
			_ = tx.Model(&model.Product{}).Where("id = ?", it.ProductID).
				Update("stock", gorm.Expr("stock + ?", it.Quantity)).Error
		}
		return nil
	})
}

type SkuStockItem struct {
	ProductID uint64
	SkuID     uint64
	Quantity  int
}

// Tags / templates CRUD
func (r *ProductAdminRepository) ListTags(shopID uint64) ([]model.ProductTag, error) {
	var list []model.ProductTag
	err := r.db.Where("shop_id IN (0, ?)", shopID).Order("id DESC").Find(&list).Error
	return list, err
}

func (r *ProductAdminRepository) SaveTag(tag *model.ProductTag) error {
	if tag.ID == 0 {
		return r.db.Create(tag).Error
	}
	return r.db.Model(tag).Updates(map[string]interface{}{"name": tag.Name, "color": tag.Color, "status": tag.Status}).Error
}

func (r *ProductAdminRepository) DeleteTag(id, shopID uint64) error {
	return r.db.Where("id = ? AND shop_id = ?", id, shopID).Delete(&model.ProductTag{}).Error
}

func (r *ProductAdminRepository) ListAttrTemplates(shopID uint64) ([]model.ProductAttrTemplate, error) {
	var list []model.ProductAttrTemplate
	err := r.db.Where("shop_id IN (0, ?)", shopID).Order("id DESC").Find(&list).Error
	return list, err
}

func (r *ProductAdminRepository) SaveAttrTemplate(t *model.ProductAttrTemplate) error {
	if t.ID == 0 {
		return r.db.Create(t).Error
	}
	return r.db.Model(t).Updates(map[string]interface{}{"name": t.Name, "attrs_json": t.AttrsJSON, "status": t.Status}).Error
}

func (r *ProductAdminRepository) DeleteAttrTemplate(id, shopID uint64) error {
	return r.db.Where("id = ? AND shop_id = ?", id, shopID).Delete(&model.ProductAttrTemplate{}).Error
}

func (r *ProductAdminRepository) CreateSchedule(s *model.ProductSchedule) error {
	_ = r.db.Model(&model.ProductSchedule{}).
		Where("product_id = ? AND action = ? AND status = ?", s.ProductID, s.Action, "pending").
		Update("status", "cancelled").Error
	return r.db.Create(s).Error
}

func (r *ProductAdminRepository) CancelSchedule(id, shopID uint64) error {
	return r.db.Model(&model.ProductSchedule{}).
		Where("id = ? AND shop_id = ? AND status = ?", id, shopID, "pending").
		Update("status", "cancelled").Error
}

func (r *ProductAdminRepository) ClaimDueSchedules(limit int) ([]model.ProductSchedule, error) {
	var list []model.ProductSchedule
	now := time.Now()
	err := r.db.Where("status = ? AND run_at <= ? AND locked_at IS NULL", "pending", now).
		Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	out := make([]model.ProductSchedule, 0, len(list))
	for _, s := range list {
		res := r.db.Model(&model.ProductSchedule{}).
			Where("id = ? AND status = ? AND locked_at IS NULL", s.ID, "pending").
			Update("locked_at", now)
		if res.RowsAffected == 1 {
			out = append(out, s)
		}
	}
	return out, nil
}

func (r *ProductAdminRepository) FinishSchedule(id uint64, ok bool) error {
	st := "done"
	if !ok {
		st = "cancelled"
	}
	return r.db.Model(&model.ProductSchedule{}).Where("id = ?", id).Update("status", st).Error
}

func (r *ProductAdminRepository) FirstSkuID(productID uint64) uint64 {
	var id uint64
	_ = r.db.Model(&model.ProductSku{}).Select("id").
		Where("product_id = ? AND deleted_at IS NULL", productID).
		Order("id ASC").Limit(1).Scan(&id).Error
	return id
}

func (r *ProductAdminRepository) GetSku(id uint64) (*model.ProductSku, error) {
	var s model.ProductSku
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&s).Error
	return &s, err
}

func (r *ProductAdminRepository) AggregatePublic(productID uint64) error {
	return r.aggregateFromSKUs(r.db, productID)
}
