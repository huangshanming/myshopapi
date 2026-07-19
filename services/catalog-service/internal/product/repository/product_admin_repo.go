package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/product/skuutil"
	"mymall/services/catalog-service/internal/product/types"

	"gorm.io/gorm"
)

type ProductAdminRepository struct {
	db *gorm.DB
}

func NewProductAdminRepository(db *gorm.DB) *ProductAdminRepository {
	return &ProductAdminRepository{db: db}
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
	Recycle       bool // true=仅回收站
	CreatedFrom   *time.Time
	CreatedTo     *time.Time
	PublishFrom   *time.Time
	PublishTo     *time.Time
	PlatformScope bool // true=允许 ShopID=0 查全站
}

func (r *ProductAdminRepository) List(f ProductListFilter) ([]model.Product, int64, error) {
	q := r.db.Model(&model.Product{})
	if f.PlatformScope {
		if f.ShopID > 0 {
			q = q.Where("shop_id = ?", f.ShopID)
		}
	} else {
		q = q.Where("shop_id = ?", f.ShopID)
	}
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
	if f.CreatedFrom != nil {
		q = q.Where("created_at >= ?", *f.CreatedFrom)
	}
	if f.CreatedTo != nil {
		q = q.Where("created_at <= ?", *f.CreatedTo)
	}
	if f.PublishFrom != nil {
		q = q.Where("publish_time >= ?", *f.PublishFrom)
	}
	if f.PublishTo != nil {
		q = q.Where("publish_time <= ?", *f.PublishTo)
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
	case "collect_desc":
		order = "collect_count DESC"
	case "collect_asc":
		order = "collect_count ASC"
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
	action := "create"
	if id > 0 {
		action = "update"
	}
	after, _ := json.Marshal(map[string]interface{}{
		"name": product.Name, "status": product.Status, "sale_price": product.SalePrice,
		"stock": product.Stock, "category_id": product.CategoryID, "product_type": product.ProductType,
	})
	_ = r.AddOpLog(shopID, &product.ID, operatorID, action, "", string(after))
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

func (r *ProductAdminRepository) GetByID(id uint64) (*model.Product, error) {
	var p model.Product
	if err := r.db.Where("id = ?", id).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductAdminRepository) SetStatusByID(id uint64, status string) error {
	p, err := r.GetByID(id)
	if err != nil {
		return errors.New("商品不存在")
	}
	return r.SetStatus(id, p.ShopID, status)
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
	// MySQL JSON 列不能写空字符串，空值必须用 null
	before = jsonOrNull(before)
	after = jsonOrNull(after)
	return r.db.Create(&model.ProductOpLog{
		ShopID: shopID, ProductID: productID, OperatorID: operatorID,
		Action: action, BeforeJSON: before, AfterJSON: after,
	}).Error
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

func (r *ProductAdminRepository) ListOpLogs(shopID uint64, productID uint64, page, pageSize int) ([]OpLogItem, int64, error) {
	q := r.db.Model(&model.ProductOpLog{}).Where("shop_id = ?", shopID)
	if productID > 0 {
		q = q.Where("product_id = ?", productID)
	}
	var total int64
	_ = q.Count(&total)
	var logs []model.ProductOpLog
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
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
			ID   uint64
			Name string
		}
		var rows []row
		_ = r.db.Table("products").Select("id, name").Where("id IN ?", productIDs).Find(&rows).Error
		for _, x := range rows {
			nameByProduct[x.ID] = x.Name
		}
	}
	nameByUser := map[uint64]string{}
	if len(operatorIDs) > 0 {
		type row struct {
			ID       uint64
			Nickname string
			Mobile   string
		}
		var rows []row
		_ = r.db.Table("users").Select("id, nickname, mobile").Where("id IN ?", operatorIDs).Find(&rows).Error
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
				// 尝试从 after_json 取 name（永久删除后商品可能已不在）
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
