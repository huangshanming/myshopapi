package merchant

import (

	"mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/svc"
)

type CatalogHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.CatalogLogic
}

func NewCatalogHandler(svcCtx *svc.ServiceContext) *CatalogHandler {
	return &CatalogHandler{
		svcCtx: svcCtx,
		logic:  logic.NewCatalogLogic(svcCtx),
	}
}

type ProductHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.ProductAdminLogic
}

func NewProductHandler(svcCtx *svc.ServiceContext) *ProductHandler {
	return &ProductHandler{
		svcCtx: svcCtx,
		logic:  logic.NewProductAdminLogic(svcCtx),
	}
}
