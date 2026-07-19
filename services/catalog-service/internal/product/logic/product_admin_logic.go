package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/product/promotion"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/product/types"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/uploadpath"
)

type ProductAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	guard  promotion.Guard
}

func NewProductAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductAdminLogic {
	return &ProductAdminLogic{ctx: ctx, svcCtx: svcCtx, guard: promotion.NewNoop()}
}

func (l *ProductAdminLogic) List(f repository.ProductListFilter) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.ProductAdmin.List(f)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"total": total, "list": list}, nil
}

func (l *ProductAdminLogic) Detail(id, shopID uint64) (map[string]interface{}, error) {
	p, skus, imgs, attrs, err := l.svcCtx.ProductAdmin.GetDetail(id, shopID)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	return map[string]interface{}{
		"product": p, "skus": skus, "images": imgs, "attrs": attrs,
	}, nil
}

func (l *ProductAdminLogic) Save(shopID, operatorID, id uint64, req types.MerchantProductSaveReq) (*model.Product, error) {
	if req.Name == "" || req.CategoryID == 0 {
		return nil, errors.New("名称与分类必填")
	}
	return l.svcCtx.ProductAdmin.SaveProduct(shopID, operatorID, id, req)
}

func (l *ProductAdminLogic) SetStatus(shopID, operatorID, id uint64, status string) error {
	switch status {
	case model.ProductOnSale, model.ProductOffSale, model.ProductDeleted, model.ProductDraft:
	default:
		return errors.New("非法状态")
	}
	if status == model.ProductOffSale || status == model.ProductDeleted {
		if !l.guard.CanOffSale(id) {
			return errors.New("商品参与活动中，不可下架/删除")
		}
	}
	if err := l.svcCtx.ProductAdmin.SetStatus(id, shopID, status); err != nil {
		return err
	}
	pid := id
	after, _ := json.Marshal(map[string]interface{}{"status": status})
	_ = l.svcCtx.ProductAdmin.AddOpLog(shopID, &pid, operatorID, "status:"+status, "", string(after))
	return nil
}

func (l *ProductAdminLogic) Copy(shopID, operatorID, id uint64) (*model.Product, error) {
	p, err := l.svcCtx.ProductAdmin.CopyProduct(id, shopID, operatorID)
	if err != nil {
		return nil, err
	}
	after, _ := json.Marshal(map[string]interface{}{"copy_from": id, "new_id": p.ID, "name": p.Name})
	_ = l.svcCtx.ProductAdmin.AddOpLog(shopID, &p.ID, operatorID, "copy", "", string(after))
	return p, nil
}

func (l *ProductAdminLogic) Batch(shopID, operatorID uint64, req types.BatchProductReq) (*model.ProductBatchJob, error) {
	if len(req.ProductIDs) == 0 {
		return nil, errors.New("未选择商品")
	}
	payload, _ := json.Marshal(req)
	job := &model.ProductBatchJob{
		ShopID: shopID, JobType: req.Action, PayloadJSON: string(payload),
		Total: len(req.ProductIDs), Status: "pending", OperatorID: operatorID,
	}
	if err := l.svcCtx.ProductAdmin.CreateBatchJob(job); err != nil {
		return nil, err
	}
	go l.runBatchJob(job.ID, shopID, operatorID, req)
	return job, nil
}

func (l *ProductAdminLogic) runBatchJob(jobID, shopID, operatorID uint64, req types.BatchProductReq) {
	db := l.svcCtx.DB
	_ = db.Model(&model.ProductBatchJob{}).Where("id = ?", jobID).Update("status", "running")
	ok, fail := 0, 0
	const batch = 50
	for i := 0; i < len(req.ProductIDs); i += batch {
		end := i + batch
		if end > len(req.ProductIDs) {
			end = len(req.ProductIDs)
		}
		chunk := req.ProductIDs[i:end]
		for _, id := range chunk {
			var err error
			switch req.Action {
			case "on_sale":
				err = l.SetStatus(shopID, operatorID, id, model.ProductOnSale)
			case "off_sale":
				err = l.SetStatus(shopID, operatorID, id, model.ProductOffSale)
			case "recycle":
				err = l.SetStatus(shopID, operatorID, id, model.ProductDeleted)
			case "category":
				err = db.Model(&model.Product{}).Where("id = ? AND shop_id = ?", id, shopID).
					Update("category_id", req.CategoryID).Error
			case "price":
				err = l.batchPrice(id, shopID, req)
			default:
				err = errors.New("未知动作")
			}
			if err != nil {
				fail++
			} else {
				ok++
			}
		}
		_ = db.Model(&model.ProductBatchJob{}).Where("id = ?", jobID).
			Updates(map[string]interface{}{"progress": ok + fail}).Error
	}
	st := "success"
	msg := fmt.Sprintf("成功%d 失败%d", ok, fail)
	if fail > 0 && ok > 0 {
		st = "partial"
	} else if fail > 0 {
		st = "failed"
	}
	_ = db.Model(&model.ProductBatchJob{}).Where("id = ?", jobID).
		Updates(map[string]interface{}{"status": st, "result_msg": msg, "progress": ok + fail}).Error
}

func (l *ProductAdminLogic) batchPrice(id, shopID uint64, req types.BatchProductReq) error {
	var skus []model.ProductSku
	if err := l.svcCtx.DB.Where("product_id = ? AND shop_id = ? AND deleted_at IS NULL", id, shopID).Find(&skus).Error; err != nil {
		return err
	}
	for _, s := range skus {
		price := s.SalePrice
		if req.PriceRate > 0 {
			price = price * req.PriceRate
		} else {
			price = price + req.PriceDelta
		}
		if price < 0.01 {
			price = 0.01
		}
		_ = l.svcCtx.DB.Model(&s).Update("sale_price", price).Error
	}
	return l.svcCtx.ProductAdmin.AggregatePublic(id)
}

func (l *ProductAdminLogic) Restore(shopID, operatorID uint64, ids []uint64) error {
	for _, id := range ids {
		if err := l.SetStatus(shopID, operatorID, id, model.ProductOffSale); err != nil {
			return err
		}
	}
	return nil
}

func (l *ProductAdminLogic) PermanentDelete(shopID, operatorID uint64, ids []uint64) error {
	for _, id := range ids {
		if !l.guard.CanDelete(id) {
			return fmt.Errorf("商品 %d 参与活动不可删除", id)
		}
	}
	if err := l.svcCtx.ProductAdmin.PermanentDelete(shopID, ids); err != nil {
		return err
	}
	after, _ := json.Marshal(map[string]interface{}{"deleted_ids": ids})
	_ = l.svcCtx.ProductAdmin.AddOpLog(shopID, nil, operatorID, "permanent_delete", "", string(after))
	return nil
}

func (l *ProductAdminLogic) AdjustStock(shopID uint64, req types.StockAdjustReq) error {
	return l.svcCtx.ProductAdmin.AdjustSkuStock(shopID, req)
}

func (l *ProductAdminLogic) BatchStock(shopID uint64, req types.BatchStockReq) error {
	for _, it := range req.Items {
		if err := l.svcCtx.ProductAdmin.AdjustSkuStock(shopID, it); err != nil {
			return err
		}
	}
	return nil
}

func (l *ProductAdminLogic) StockWarnings(shopID uint64, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.ProductAdmin.ListStockWarnings(shopID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"total": total, "list": list}, nil
}

func (l *ProductAdminLogic) SaveUpload(shopID uint64, filename string, data []byte) (string, error) {
	if len(data) > 5*1024*1024 {
		return "", errors.New("文件不能超过5MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return "", errors.New("仅支持图片")
	}
	dir := uploadpath.Abs("products", fmt.Sprintf("%d", shopID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return "/uploads/products/" + fmt.Sprintf("%d", shopID) + "/" + name, nil
}

func (l *ProductAdminLogic) CreateSchedule(shopID, operatorID, productID uint64, req types.ScheduleReq) error {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", req.RunAt, time.Local)
	if err != nil {
		return errors.New("时间格式应为 2006-01-02 15:04:05")
	}
	if req.Action != model.ProductOnSale && req.Action != model.ProductOffSale {
		return errors.New("action 无效")
	}
	s := &model.ProductSchedule{
		ProductID: productID, ShopID: shopID, Action: req.Action,
		RunAt: common.LocalTime(t), Status: "pending",
	}
	if err := l.svcCtx.ProductAdmin.CreateSchedule(s); err != nil {
		return err
	}
	after, _ := json.Marshal(map[string]interface{}{"action": req.Action, "run_at": req.RunAt})
	_ = l.svcCtx.ProductAdmin.AddOpLog(shopID, &productID, operatorID, "schedule", "", string(after))
	return nil
}

func (l *ProductAdminLogic) RunSchedules() {
	list, err := l.svcCtx.ProductAdmin.ClaimDueSchedules(20)
	if err != nil {
		return
	}
	for _, s := range list {
		err := l.svcCtx.ProductAdmin.SetStatus(s.ProductID, s.ShopID, s.Action)
		_ = l.svcCtx.ProductAdmin.FinishSchedule(s.ID, err == nil)
	}
}

func (l *ProductAdminLogic) OpLogs(shopID, productID uint64, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.ProductAdmin.ListOpLogs(shopID, productID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"total": total, "list": list}, nil
}

func (l *ProductAdminLogic) Job(shopID, id uint64) (*model.ProductBatchJob, error) {
	return l.svcCtx.ProductAdmin.GetBatchJob(id, shopID)
}

// 简易导出 CSV
func (l *ProductAdminLogic) ExportCSV(shopID uint64) (string, error) {
	list, _, err := l.svcCtx.ProductAdmin.List(repository.ProductListFilter{ShopID: shopID, Page: 1, PageSize: 5000})
	if err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("id,product_no,name,status,sale_price,stock,category_id,product_type\n")
	for _, p := range list {
		b.WriteString(fmt.Sprintf("%d,%s,%s,%s,%.2f,%d,%d,%s\n",
			p.ID, p.ProductNo, escapeCSV(p.Name), p.Status, p.SalePrice, p.Stock, p.CategoryID, p.ProductType))
	}
	dir := uploadpath.Abs("exports", fmt.Sprintf("%d", shopID))
	_ = os.MkdirAll(dir, 0o755)
	name := fmt.Sprintf("products_%d.csv", time.Now().Unix())
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return "", err
	}
	return "/uploads/exports/" + fmt.Sprintf("%d", shopID) + "/" + name, nil
}

func escapeCSV(s string) string {
	s = strings.ReplaceAll(s, "\"", "\"\"")
	if strings.ContainsAny(s, ",\"\n") {
		return "\"" + s + "\""
	}
	return s
}

// ImportCSV 简化：每行 name,category_id,sale_price,stock
func (l *ProductAdminLogic) ImportCSV(shopID, operatorID uint64, content string) (map[string]interface{}, error) {
	lines := strings.Split(content, "\n")
	ok, fail := 0, 0
	errs := []string{}
	seen := map[string]struct{}{}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || i == 0 && strings.HasPrefix(strings.ToLower(line), "name") {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 4 {
			fail++
			errs = append(errs, fmt.Sprintf("行%d格式错误", i+1))
			continue
		}
		name := strings.TrimSpace(parts[0])
		if name == "" {
			fail++
			continue
		}
		if _, dup := seen[name]; dup {
			fail++
			errs = append(errs, fmt.Sprintf("行%d重复名称%s", i+1, name))
			continue
		}
		seen[name] = struct{}{}
		var catID uint64
		var price float64
		var stock int
		fmt.Sscanf(parts[1], "%d", &catID)
		fmt.Sscanf(parts[2], "%f", &price)
		fmt.Sscanf(parts[3], "%d", &stock)
		if catID == 0 || price <= 0 {
			fail++
			errs = append(errs, fmt.Sprintf("行%d价格或分类无效", i+1))
			continue
		}
		_, err := l.Save(shopID, operatorID, 0, types.MerchantProductSaveReq{
			Name: name, CategoryID: catID, SalePrice: price, Stock: stock,
			Status: model.ProductDraft, ProductType: model.ProductTypePhysical,
			Skus: []types.SkuInput{{SalePrice: price, Stock: stock, SpecValues: map[string]string{}, Status: model.SKUEnabled}},
		})
		if err != nil {
			fail++
			errs = append(errs, err.Error())
		} else {
			ok++
		}
	}
	return map[string]interface{}{"ok": ok, "fail": fail, "errors": errs}, nil
}
