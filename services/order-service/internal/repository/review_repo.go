package repository

import (
	"context"
	"time"

	"mymall/common"
	"mymall/services/order-service/internal/model"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, rev *model.ProductReview, images []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rev).Error; err != nil {
			return err
		}
		for i, url := range images {
			if url == "" {
				continue
			}
			img := model.ProductReviewImage{ReviewID: rev.ID, URL: url, Sort: i}
			if err := tx.Create(&img).Error; err != nil {
				return err
			}
		}
		now := common.LocalTime(time.Now())
		res := tx.Model(&model.Order{}).
			Where("id = ? AND user_id = ? AND status = ?", rev.OrderID, rev.UserID, model.OrderStatusCompleted).
			Updates(map[string]interface{}{
				"status":      model.OrderStatusReviewed,
				"reviewed_at": &now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *ReviewRepository) GetByOrderID(ctx context.Context, orderID uint64) (*model.ProductReview, error) {
	var rev model.ProductReview
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&rev).Error; err != nil {
		return nil, err
	}
	var imgs []model.ProductReviewImage
	_ = r.db.WithContext(ctx).Where("review_id = ?", rev.ID).Order("sort ASC, id ASC").Find(&imgs).Error
	rev.Images = imgs
	return &rev, nil
}

func (r *ReviewRepository) ExistsByOrderID(ctx context.Context, orderID uint64) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.ProductReview{}).Where("order_id = ?", orderID).Count(&n).Error
	return n > 0, err
}

func (r *ReviewRepository) ListByProduct(ctx context.Context, productID uint64, page, pageSize int) ([]model.ProductReview, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	q := r.db.WithContext(ctx).Model(&model.ProductReview{}).
		Where("product_id = ? AND status = ?", productID, model.ReviewStatusVisible)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.ProductReview
	if err := q.Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	r.attachImages(ctx, list)
	return list, total, nil
}

func (r *ReviewRepository) ListMerchant(ctx context.Context, shopID uint64, ratingLevel string, page, pageSize int) ([]model.ProductReview, int64, error) {
	return r.listFiltered(ctx, shopID, 0, ratingLevel, page, pageSize)
}

func (r *ReviewRepository) ListAdmin(ctx context.Context, shopID uint64, ratingLevel string, page, pageSize int) ([]model.ProductReview, int64, error) {
	return r.listFiltered(ctx, shopID, 0, ratingLevel, page, pageSize)
}

func (r *ReviewRepository) listFiltered(ctx context.Context, shopID, productID uint64, ratingLevel string, page, pageSize int) ([]model.ProductReview, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.WithContext(ctx).Model(&model.ProductReview{}).Where("status = ?", model.ReviewStatusVisible)
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	if productID > 0 {
		q = q.Where("product_id = ?", productID)
	}
	switch ratingLevel {
	case "good":
		q = q.Where("rating >= 4")
	case "mid":
		q = q.Where("rating = 3")
	case "bad":
		q = q.Where("rating <= 2")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.ProductReview
	if err := q.Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	r.attachImages(ctx, list)
	return list, total, nil
}

func (r *ReviewRepository) attachImages(ctx context.Context, list []model.ProductReview) {
	if len(list) == 0 {
		return
	}
	ids := make([]uint64, len(list))
	for i, v := range list {
		ids[i] = v.ID
	}
	var imgs []model.ProductReviewImage
	_ = r.db.WithContext(ctx).Where("review_id IN ?", ids).Order("sort ASC, id ASC").Find(&imgs).Error
	m := map[uint64][]model.ProductReviewImage{}
	for _, img := range imgs {
		m[img.ReviewID] = append(m[img.ReviewID], img)
	}
	for i := range list {
		list[i].Images = m[list[i].ID]
	}
}

func (r *ReviewRepository) GetByID(ctx context.Context, id uint64) (*model.ProductReview, error) {
	var rev model.ProductReview
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rev).Error; err != nil {
		return nil, err
	}
	var imgs []model.ProductReviewImage
	_ = r.db.WithContext(ctx).Where("review_id = ?", rev.ID).Order("sort ASC, id ASC").Find(&imgs).Error
	rev.Images = imgs
	return &rev, nil
}

func (r *ReviewRepository) Reply(ctx context.Context, id, shopID uint64, reply string) error {
	now := common.LocalTime(time.Now())
	res := r.db.WithContext(ctx).Model(&model.ProductReview{}).
		Where("id = ? AND shop_id = ? AND status = ?", id, shopID, model.ReviewStatusVisible).
		Updates(map[string]interface{}{
			"merchant_reply": reply,
			"replied_at":     &now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ReviewRepository) SoftDelete(ctx context.Context, id uint64, shopID uint64) error {
	q := r.db.WithContext(ctx).Model(&model.ProductReview{}).Where("id = ? AND status = ?", id, model.ReviewStatusVisible)
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	res := q.Update("status", model.ReviewStatusDeleted)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ReviewRepository) ProductStats(ctx context.Context, productID uint64) (avg float64, count int64, goodRate float64, err error) {
	type row struct {
		Cnt  int64
		Avg  float64
		Good int64
	}
	var s row
	err = r.db.WithContext(ctx).Model(&model.ProductReview{}).
		Select("COUNT(*) AS cnt, COALESCE(AVG(rating),0) AS avg, SUM(CASE WHEN rating >= 4 THEN 1 ELSE 0 END) AS good").
		Where("product_id = ? AND status = ?", productID, model.ReviewStatusVisible).
		Scan(&s).Error
	if err != nil {
		return 0, 0, 0, err
	}
	count = s.Cnt
	avg = s.Avg
	if count > 0 {
		goodRate = float64(s.Good) * 100 / float64(count)
	}
	return avg, count, goodRate, nil
}

func (r *ReviewRepository) UpdateProductStats(ctx context.Context, productID uint64, avg float64, count int64, goodRate float64) error {
	return r.db.WithContext(ctx).Exec(
		"UPDATE products SET avg_rating = ?, review_count = ?, good_rate = ? WHERE id = ?",
		avg, count, goodRate, productID,
	).Error
}
