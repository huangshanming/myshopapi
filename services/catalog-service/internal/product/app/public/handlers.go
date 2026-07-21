package public

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
