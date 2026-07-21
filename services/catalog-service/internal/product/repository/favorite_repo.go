package repository

import (
	"context"
	"errors"

	"mymall/services/catalog-service/internal/product/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Add(ctx context.Context, userID, productID uint64) (created bool, err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var p model.Product
		if err := tx.Select("id").Where("id = ?", productID).First(&p).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("商品不存在")
			}
			return err
		}
		fav := model.ProductFavorite{UserID: userID, ProductID: productID}
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&fav)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected > 0 {
			created = true
			return tx.Model(&model.Product{}).Where("id = ?", productID).
				UpdateColumn("collect_count", gorm.Expr("collect_count + 1")).Error
		}
		return nil
	})
	return created, err
}

func (r *FavoriteRepository) Remove(ctx context.Context, userID, productID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&model.ProductFavorite{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected > 0 {
			return tx.Model(&model.Product{}).Where("id = ?", productID).
				UpdateColumn("collect_count", gorm.Expr("GREATEST(0, collect_count - 1)")).Error
		}
		return nil
	})
}

func (r *FavoriteRepository) RemoveBatch(ctx context.Context, userID uint64, productIDs []uint64) error {
	if len(productIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, pid := range productIDs {
			res := tx.Where("user_id = ? AND product_id = ?", userID, pid).Delete(&model.ProductFavorite{})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected > 0 {
				if err := tx.Model(&model.Product{}).Where("id = ?", pid).
					UpdateColumn("collect_count", gorm.Expr("GREATEST(0, collect_count - 1)")).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *FavoriteRepository) IsFavorited(ctx context.Context, userID, productID uint64) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.ProductFavorite{}).
		Where("user_id = ? AND product_id = ?", userID, productID).Count(&n).Error
	return n > 0, err
}

func (r *FavoriteRepository) CountByProduct(ctx context.Context, productID uint64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.ProductFavorite{}).Where("product_id = ?", productID).Count(&n).Error
	return n, err
}

func (r *FavoriteRepository) List(ctx context.Context, userID uint64, page, pageSize int) ([]model.FavoriteListItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.ProductFavorite{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var favs []model.ProductFavorite
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&favs).Error; err != nil {
		return nil, 0, err
	}
	if len(favs) == 0 {
		return []model.FavoriteListItem{}, total, nil
	}
	ids := make([]uint64, len(favs))
	for i, f := range favs {
		ids[i] = f.ProductID
	}
	var products []model.Product
	_ = r.db.WithContext(ctx).Select("id, name, main_image, sale_price, status, collect_count").
		Where("id IN ?", ids).Find(&products).Error
	pmap := make(map[uint64]model.Product, len(products))
	for _, p := range products {
		pmap[p.ID] = p
	}
	out := make([]model.FavoriteListItem, 0, len(favs))
	for _, f := range favs {
		p, ok := pmap[f.ProductID]
		item := model.FavoriteListItem{
			ID:        f.ID,
			ProductID: f.ProductID,
			CreatedAt: f.CreatedAt,
			Invalid:   true,
		}
		if ok {
			item.Name = p.Name
			item.MainImage = p.MainImage
			item.SalePrice = p.SalePrice
			item.Status = p.Status
			item.CollectCount = p.CollectCount
			item.Invalid = p.Status != "on_sale"
		}
		out = append(out, item)
	}
	return out, total, nil
}

func (r *FavoriteRepository) UpdateReviewStats(ctx context.Context, productID uint64, avg float64, count int, goodRate float64) error {
	return r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
		"avg_rating":   avg,
		"review_count": count,
		"good_rate":    goodRate,
	}).Error
}
