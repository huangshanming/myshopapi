package admin

import (
	"context"

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
		logic:  logic.NewCatalogLogic(context.Background(), svcCtx),
	}
}

type FavoriteHandler struct {
	logic *logic.FavoriteLogic
}

func NewFavoriteHandler(svcCtx *svc.ServiceContext) *FavoriteHandler {
	return &FavoriteHandler{logic: logic.NewFavoriteLogic(context.Background(), svcCtx)}
}

type PlatformProductHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.PlatformProductLogic
}

func NewPlatformProductHandler(svcCtx *svc.ServiceContext) *PlatformProductHandler {
	return &PlatformProductHandler{
		svcCtx: svcCtx,
		logic:  logic.NewPlatformProductLogic(context.Background(), svcCtx),
	}
}

type ShopUploadHandler struct{}

func NewShopUploadHandler() *ShopUploadHandler {
	return &ShopUploadHandler{}
}
