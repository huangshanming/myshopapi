package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mymall/pkg/pagination"
	"mymall/services/catalog-service/internal/model"
	"mymall/services/catalog-service/internal/repository"

	"github.com/redis/go-redis/v9"
)

const productListCacheTTL = 5 * time.Minute

type CatalogService struct {
	products   *repository.ProductRepository
	categories *repository.CategoryRepository
	redis      *redis.Client
}

func NewCatalogService(products *repository.ProductRepository, categories *repository.CategoryRepository, redisClient *redis.Client) *CatalogService {
	return &CatalogService{products: products, categories: categories, redis: redisClient}
}

func (s *CatalogService) productListCacheKey(page *pagination.PageReq) string {
	return fmt.Sprintf("catalog:products:list:%d:%d", page.Page, page.PageSize)
}

func (s *CatalogService) GetProductList(page *pagination.PageReq) (map[string]interface{}, error) {
	return s.GetProductListFiltered(page, 0, "on_sale")
}

func (s *CatalogService) GetProductListFiltered(page *pagination.PageReq, shopID uint64, status string) (map[string]interface{}, error) {
	if s.redis != nil && shopID == 0 && status == "on_sale" {
		key := s.productListCacheKey(page)
		cached, err := s.redis.Get(context.Background(), key).Bytes()
		if err == nil {
			var res map[string]interface{}
			if json.Unmarshal(cached, &res) == nil {
				return res, nil
			}
		}
	}

	res, err := s.products.GetListByShop(page, shopID, status)
	if err != nil {
		return res, err
	}

	if s.redis != nil && shopID == 0 && status == "on_sale" {
		if data, err := json.Marshal(res); err == nil {
			_ = s.redis.Set(context.Background(), s.productListCacheKey(page), data, productListCacheTTL).Err()
		}
	}
	return res, nil
}

func (s *CatalogService) CreateProduct(p *model.Product) error {
	return s.products.Create(p)
}

func (s *CatalogService) UpdateProductByShop(id, shopID uint64, updates map[string]interface{}) error {
	return s.products.UpdateByShop(id, shopID, updates)
}

func (s *CatalogService) ForceOffSale(id uint64) error {
	return s.products.ForceOffSale(id)
}

func (s *CatalogService) CreateCategory(c *model.ProductCategory) error {
	return s.categories.Create(c)
}

func (s *CatalogService) UpdateCategory(id uint64, updates map[string]interface{}) error {
	return s.categories.Update(id, updates)
}

func (s *CatalogService) DeleteCategory(id uint64) error {
	return s.categories.Delete(id)
}

func (s *CatalogService) GetProductDetail(id uint64) (*model.Product, error) {
	return s.products.GetDetail(id)
}

func (s *CatalogService) GetCategoryList(page *pagination.PageReq) (*pagination.PageRes[model.ProductCategory], error) {
	return s.categories.GetList(page)
}

func (s *CatalogService) GetCategoryDetail(id uint64) (*model.ProductCategory, error) {
	return s.categories.GetDetail(id)
}

func (s *CatalogService) BatchGetProducts(ids []uint64) ([]model.Product, error) {
	return s.products.BatchGetByIDs(ids)
}

func (s *CatalogService) ReserveStock(items []repository.StockItem) error {
	return s.products.ReserveStock(items)
}

func (s *CatalogService) ReleaseStock(items []repository.StockItem) error {
	return s.products.ReleaseStock(items)
}
