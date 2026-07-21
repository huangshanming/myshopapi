package admin

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

type FavoriteHandler struct {
	logic *logic.FavoriteLogic
}

func NewFavoriteHandler(svcCtx *svc.ServiceContext) *FavoriteHandler {
	return &FavoriteHandler{logic: logic.NewFavoriteLogic(svcCtx)}
}

type PlatformProductHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.PlatformProductLogic
}

func NewPlatformProductHandler(svcCtx *svc.ServiceContext) *PlatformProductHandler {
	return &PlatformProductHandler{
		svcCtx: svcCtx,
		logic:  logic.NewPlatformProductLogic(svcCtx),
	}
}

type ShopUploadHandler struct{}

func NewShopUploadHandler() *ShopUploadHandler {
	return &ShopUploadHandler{}
}
