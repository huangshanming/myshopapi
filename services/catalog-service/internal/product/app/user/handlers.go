package user

import (
	"mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/svc"
)

type FavoriteHandler struct {
	logic *logic.FavoriteLogic
}

func NewFavoriteHandler(svcCtx *svc.ServiceContext) *FavoriteHandler {
	return &FavoriteHandler{logic: logic.NewFavoriteLogic(svcCtx)}
}
