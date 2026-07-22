package repository

import (
	"context"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/order-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	reviewColumns = "id, order_id, order_no, user_id, shop_id, product_id, order_item_id, sku_id, sku_snapshot, rating, content, is_anonymous, status, merchant_reply, replied_at, created_at, updated_at"
	reviewImageColumns = "id, review_id, url, sort"
)

type ReviewRepository struct {
	conn sqlx.SqlConn
}

func NewReviewRepository(conn sqlx.SqlConn) *ReviewRepository {
	return &ReviewRepository{conn: conn}
}

func (r *ReviewRepository) Create(ctx context.Context, rev *model.ProductReview, images []string) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		anon := 0
		if rev.IsAnonymous {
			anon = 1
		}
		id, err := lastInsertID(ctx, session,
			`INSERT INTO product_reviews (order_id, order_no, user_id, shop_id, product_id, order_item_id, sku_id, sku_snapshot, rating, content, is_anonymous, status, merchant_reply)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			rev.OrderID, rev.OrderNo, rev.UserID, rev.ShopID, rev.ProductID, rev.OrderItemID, rev.SkuID, rev.SkuSnapshot,
			rev.Rating, rev.Content, anon, rev.Status, rev.MerchantReply,
		)
		if err != nil {
			return err
		}
		rev.ID = id
		for i, url := range images {
			if url == "" {
				continue
			}
			if _, err := session.ExecCtx(ctx,
				"INSERT INTO product_review_images (review_id, url, sort) VALUES (?, ?, ?)",
				rev.ID, url, i,
			); err != nil {
				return err
			}
		}
		now := common.LocalTime(time.Now())
		res, err := session.ExecCtx(ctx,
			"UPDATE orders SET status=?, reviewed_at=? WHERE id=? AND user_id=? AND status=?",
			model.OrderStatusReviewed, &now, rev.OrderID, rev.UserID, model.OrderStatusCompleted,
		)
		if err != nil {
			return err
		}
		n, _ := res.RowsAffected()
		if n == 0 {
			return sqlx.ErrNotFound
		}
		return nil
	})
}

func (r *ReviewRepository) GetByOrderID(ctx context.Context, orderID uint64) (*model.ProductReview, error) {
	var rev model.ProductReview
	if err := r.conn.QueryRowCtx(ctx, &rev, "SELECT "+reviewColumns+" FROM product_reviews WHERE order_id=? LIMIT 1", orderID); err != nil {
		return nil, err
	}
	var imgs []model.ProductReviewImage
	_ = r.conn.QueryRowsCtx(ctx, &imgs,
		"SELECT "+reviewImageColumns+" FROM product_review_images WHERE review_id=? ORDER BY sort ASC, id ASC", rev.ID,
	)
	rev.Images = imgs
	return &rev, nil
}

func (r *ReviewRepository) ExistsByOrderID(ctx context.Context, orderID uint64) (bool, error) {
	n, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM product_reviews WHERE order_id=?", orderID)
	return n > 0, err
}

func (r *ReviewRepository) ListByProduct(ctx context.Context, productID uint64, page, pageSize int) ([]model.ProductReview, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	total, err := countCtx(ctx, r.conn,
		"SELECT COUNT(*) FROM product_reviews WHERE product_id=? AND status=?",
		productID, model.ReviewStatusVisible,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.ProductReview
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+reviewColumns+" FROM product_reviews WHERE product_id=? AND status=? ORDER BY created_at DESC, id DESC LIMIT ?, ?",
		productID, model.ReviewStatusVisible, (page-1)*pageSize, pageSize,
	)
	if err != nil {
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
	where := []string{"status=?"}
	args := []any{model.ReviewStatusVisible}
	if shopID > 0 {
		where = append(where, "shop_id=?")
		args = append(args, shopID)
	}
	if productID > 0 {
		where = append(where, "product_id=?")
		args = append(args, productID)
	}
	switch ratingLevel {
	case "good":
		where = append(where, "rating >= 4")
	case "mid":
		where = append(where, "rating = 3")
	case "bad":
		where = append(where, "rating <= 2")
	}
	w := strings.Join(where, " AND ")
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM product_reviews WHERE "+w, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), (page-1)*pageSize, pageSize)
	var list []model.ProductReview
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+reviewColumns+" FROM product_reviews WHERE "+w+" ORDER BY created_at DESC, id DESC LIMIT ?, ?",
		listArgs...,
	)
	if err != nil {
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
	_ = r.conn.QueryRowsCtx(ctx, &imgs,
		"SELECT "+reviewImageColumns+" FROM product_review_images WHERE review_id IN ("+placeholders(len(ids))+") ORDER BY sort ASC, id ASC",
		inArgs(ids)...,
	)
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
	if err := r.conn.QueryRowCtx(ctx, &rev, "SELECT "+reviewColumns+" FROM product_reviews WHERE id=? LIMIT 1", id); err != nil {
		return nil, err
	}
	var imgs []model.ProductReviewImage
	_ = r.conn.QueryRowsCtx(ctx, &imgs,
		"SELECT "+reviewImageColumns+" FROM product_review_images WHERE review_id=? ORDER BY sort ASC, id ASC", rev.ID,
	)
	rev.Images = imgs
	return &rev, nil
}

func (r *ReviewRepository) Reply(ctx context.Context, id, shopID uint64, reply string) error {
	now := common.LocalTime(time.Now())
	n, err := execAffected(ctx, r.conn,
		"UPDATE product_reviews SET merchant_reply=?, replied_at=? WHERE id=? AND shop_id=? AND status=?",
		reply, &now, id, shopID, model.ReviewStatusVisible,
	)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *ReviewRepository) SoftDelete(ctx context.Context, id uint64, shopID uint64) error {
	q := "UPDATE product_reviews SET status=? WHERE id=? AND status=?"
	args := []any{model.ReviewStatusDeleted, id, model.ReviewStatusVisible}
	if shopID > 0 {
		q += " AND shop_id=?"
		args = append(args, shopID)
	}
	n, err := execAffected(ctx, r.conn, q, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *ReviewRepository) ProductStats(ctx context.Context, productID uint64) (avg float64, count int64, goodRate float64, err error) {
	type row struct {
		Cnt  int64   `db:"cnt"`
		Avg  float64 `db:"avg"`
		Good int64   `db:"good"`
	}
	var s row
	err = r.conn.QueryRowCtx(ctx, &s,
		`SELECT COUNT(*) AS cnt, COALESCE(AVG(rating),0) AS avg, SUM(CASE WHEN rating >= 4 THEN 1 ELSE 0 END) AS good
		 FROM product_reviews WHERE product_id=? AND status=?`,
		productID, model.ReviewStatusVisible,
	)
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
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE products SET avg_rating = ?, review_count = ?, good_rate = ? WHERE id = ?",
		avg, count, goodRate, productID,
	)
	return err
}
