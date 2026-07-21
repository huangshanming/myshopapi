package logic

import (
	"context"

	"mymall/services/catalog-service/internal/client/userhttp"
	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/svc"
)

type FavoriteLogic struct {
	svcCtx *svc.ServiceContext
}

func NewFavoriteLogic(svcCtx *svc.ServiceContext) *FavoriteLogic {
	return &FavoriteLogic{svcCtx: svcCtx}
}

func (l *FavoriteLogic) Add(ctx context.Context, userID, productID uint64) error {
	created, err := l.svcCtx.Favorites.Add(ctx, userID, productID)
	if err != nil {
		return err
	}
	if created && l.svcCtx.UserHTTP != nil {
		_ = l.svcCtx.UserHTTP.TaskEvent(ctx, userhttp.TaskEventReq{
			UserID: userID, TaskCode: "first_favorite_product", Delta: 1,
			RefType: "product", RefID: productID,
		})
	}
	return nil
}

func (l *FavoriteLogic) Remove(ctx context.Context, userID, productID uint64) error {
	return l.svcCtx.Favorites.Remove(ctx, userID, productID)
}

func (l *FavoriteLogic) RemoveBatch(ctx context.Context, userID uint64, productIDs []uint64) error {
	return l.svcCtx.Favorites.RemoveBatch(ctx, userID, productIDs)
}

func (l *FavoriteLogic) List(ctx context.Context, userID uint64, page, pageSize int) ([]model.FavoriteListItem, int64, error) {
	return l.svcCtx.Favorites.List(ctx, userID, page, pageSize)
}

func (l *FavoriteLogic) IsFavorited(ctx context.Context, userID, productID uint64) (bool, error) {
	return l.svcCtx.Favorites.IsFavorited(ctx, userID, productID)
}

func (l *FavoriteLogic) FavoriteCount(productID uint64) (int64, error) {
	var p model.Product
	if err := l.svcCtx.DB.Select("collect_count").Where("id = ?", productID).First(&p).Error; err != nil {
		return 0, err
	}
	return int64(p.CollectCount), nil
}
