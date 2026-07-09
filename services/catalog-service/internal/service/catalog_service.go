package service

import (
	"mymall/pkg/pagination"
	"mymall/services/catalog-service/internal/model"
	"mymall/services/catalog-service/internal/repository"
)

type CatalogService struct {
	products   *repository.ProductRepository
	categories *repository.CategoryRepository
}

func NewCatalogService(products *repository.ProductRepository, categories *repository.CategoryRepository) *CatalogService {
	return &CatalogService{products: products, categories: categories}
}

func (s *CatalogService) GetProductList(page *pagination.PageReq) (map[string]interface{}, error) {
	return s.products.GetList(page)
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
