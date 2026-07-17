package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mymall/pkg/pagination"
	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/product/types"
	"mymall/services/catalog-service/internal/svc"
)

const productListCacheTTL = 5 * time.Minute

type CatalogLogic struct {
	svcCtx *svc.ServiceContext
}

func NewCatalogLogic(svcCtx *svc.ServiceContext) *CatalogLogic {
	return &CatalogLogic{svcCtx: svcCtx}
}

func (l *CatalogLogic) productListCacheKey(page *pagination.PageReq) string {
	return fmt.Sprintf("catalog:products:list:%d:%d", page.Page, page.PageSize)
}

func (l *CatalogLogic) GetProductList(page *pagination.PageReq) (map[string]interface{}, error) {
	return l.GetProductListFiltered(page, 0, "on_sale")
}

func (l *CatalogLogic) GetProductListFiltered(page *pagination.PageReq, shopID uint64, status string) (map[string]interface{}, error) {
	if l.svcCtx.Redis != nil && shopID == 0 && status == "on_sale" {
		key := l.productListCacheKey(page)
		cached, err := l.svcCtx.Redis.Get(context.Background(), key).Bytes()
		if err == nil {
			var res map[string]interface{}
			if json.Unmarshal(cached, &res) == nil {
				return res, nil
			}
		}
	}

	res, err := l.svcCtx.Products.GetListByShop(page, shopID, status)
	if err != nil {
		return res, err
	}

	if l.svcCtx.Redis != nil && shopID == 0 && status == "on_sale" {
		if data, err := json.Marshal(res); err == nil {
			_ = l.svcCtx.Redis.Set(context.Background(), l.productListCacheKey(page), data, productListCacheTTL).Err()
		}
	}
	return res, nil
}

func (l *CatalogLogic) CreateProduct(shopID uint64, req types.MerchantProductReq) (*model.Product, error) {
	status := req.Status
	if status == "" {
		status = "on_sale"
	}
	pet := req.PetType
	if pet == "" {
		pet = "both"
	}
	p := &model.Product{
		ShopID:     shopID,
		ProductNo:  fmt.Sprintf("P%d", time.Now().UnixNano()%1e12),
		Name:       req.Name,
		SalePrice:  req.SalePrice,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
		Subtitle:   req.Subtitle,
		MainImage:  req.MainImage,
		Status:     status,
		PetType:    pet,
		Discount:   100,
	}
	if err := l.svcCtx.Products.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (l *CatalogLogic) UpdateProductByShop(id, shopID uint64, req types.MerchantProductReq) error {
	updates := map[string]interface{}{
		"name":        req.Name,
		"sale_price":  req.SalePrice,
		"stock":       req.Stock,
		"category_id": req.CategoryID,
		"subtitle":    req.Subtitle,
		"main_image":  req.MainImage,
	}
	if req.Status != "" {
		updates["status"] = req.Status
		if req.Status == "on_sale" {
			updates["publish_time"] = time.Now()
		}
	}
	if req.PetType != "" {
		updates["pet_type"] = req.PetType
	}
	return l.svcCtx.Products.UpdateByShop(id, shopID, updates)
}

func (l *CatalogLogic) SetProductStatus(id, shopID uint64, status string) error {
	return l.svcCtx.Products.UpdateByShop(id, shopID, map[string]interface{}{"status": status})
}

func (l *CatalogLogic) ForceOffSale(id uint64) error {
	return l.svcCtx.Products.ForceOffSale(id)
}

func (l *CatalogLogic) CreateCategory(req types.CategoryReq) (*model.ProductCategory, error) {
	level := req.Level
	if level == 0 {
		level = 1
	}
	show := true
	if req.IsShow != nil {
		show = *req.IsShow
	}
	cat := &model.ProductCategory{
		ParentId:    req.ParentId,
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Level:       level,
		IsShow:      show,
	}
	if err := l.svcCtx.Categories.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (l *CatalogLogic) UpdateCategory(id uint64, req types.CategoryReq) error {
	updates := map[string]interface{}{
		"parent_id":   req.ParentId,
		"name":        req.Name,
		"icon":        req.Icon,
		"description": req.Description,
		"sort_order":  req.SortOrder,
	}
	if req.Level > 0 {
		updates["level"] = req.Level
	}
	if req.IsShow != nil {
		updates["is_show"] = *req.IsShow
	}
	return l.svcCtx.Categories.Update(id, updates)
}

func (l *CatalogLogic) DeleteCategory(id uint64) error {
	return l.svcCtx.Categories.Delete(id)
}

func (l *CatalogLogic) GetProductDetail(id uint64) (*model.Product, error) {
	return l.svcCtx.Products.GetDetail(id)
}

func (l *CatalogLogic) GetCategoryList(page *pagination.PageReq) (*pagination.PageRes[model.ProductCategory], error) {
	return l.svcCtx.Categories.GetList(page)
}

func (l *CatalogLogic) ListAllCategories() ([]model.ProductCategory, error) {
	return l.svcCtx.Categories.ListAll()
}

func (l *CatalogLogic) GetCategoryDetail(id uint64) (*model.ProductCategory, error) {
	return l.svcCtx.Categories.GetDetail(id)
}

func (l *CatalogLogic) BatchGetProducts(ids []uint64) ([]model.Product, error) {
	return l.svcCtx.Products.BatchGetByIDs(ids)
}

func (l *CatalogLogic) DefaultSkuID(productID uint64) uint64 {
	if l.svcCtx.ProductAdmin == nil {
		return 0
	}
	return l.svcCtx.ProductAdmin.FirstSkuID(productID)
}

func (l *CatalogLogic) GetSkuSpecSnapshot(skuID uint64) string {
	if skuID == 0 || l.svcCtx.ProductAdmin == nil {
		return "{}"
	}
	sku, err := l.svcCtx.ProductAdmin.GetSku(skuID)
	if err != nil || sku == nil {
		return "{}"
	}
	if sku.SpecValues == "" {
		return "{}"
	}
	return sku.SpecValues
}

func (l *CatalogLogic) ReserveStock(items []repository.StockItem) error {
	return l.svcCtx.Products.ReserveStock(items)
}

func (l *CatalogLogic) ReleaseStock(items []repository.StockItem) error {
	return l.svcCtx.Products.ReleaseStock(items)
}
