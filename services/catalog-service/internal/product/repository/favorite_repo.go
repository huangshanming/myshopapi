package repository

import (
	"context"
	"errors"

	"mymall/services/catalog-service/internal/product/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const productFavoriteColumns = "id, user_id, product_id, created_at"

type FavoriteRepository struct {
	conn sqlx.SqlConn
}

func NewFavoriteRepository(conn sqlx.SqlConn) *FavoriteRepository {
	return &FavoriteRepository{conn: conn}
}

func (r *FavoriteRepository) Add(ctx context.Context, userID, productID uint64) (created bool, err error) {
	err = r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var id uint64
		if err := session.QueryRowPartialCtx(ctx, &id, "SELECT id FROM products WHERE id=? LIMIT 1", productID); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return errors.New("商品不存在")
			}
			return err
		}
		res, err := session.ExecCtx(ctx,
			"INSERT IGNORE INTO product_favorites (user_id, product_id) VALUES (?, ?)", userID, productID)
		if err != nil {
			return err
		}
		n, _ := res.RowsAffected()
		if n > 0 {
			created = true
			_, err = session.ExecCtx(ctx,
				"UPDATE products SET collect_count=collect_count+1 WHERE id=?", productID)
			return err
		}
		return nil
	})
	return created, err
}

func (r *FavoriteRepository) Remove(ctx context.Context, userID, productID uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		n, err := execSession(ctx, session,
			"DELETE FROM product_favorites WHERE user_id=? AND product_id=?", userID, productID)
		if err != nil {
			return err
		}
		if n > 0 {
			_, err = session.ExecCtx(ctx,
				"UPDATE products SET collect_count=GREATEST(0, collect_count-1) WHERE id=?", productID)
			return err
		}
		return nil
	})
}

func (r *FavoriteRepository) RemoveBatch(ctx context.Context, userID uint64, productIDs []uint64) error {
	if len(productIDs) == 0 {
		return nil
	}
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, pid := range productIDs {
			n, err := execSession(ctx, session,
				"DELETE FROM product_favorites WHERE user_id=? AND product_id=?", userID, pid)
			if err != nil {
				return err
			}
			if n > 0 {
				if _, err := session.ExecCtx(ctx,
					"UPDATE products SET collect_count=GREATEST(0, collect_count-1) WHERE id=?", pid); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *FavoriteRepository) IsFavorited(ctx context.Context, userID, productID uint64) (bool, error) {
	n, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM product_favorites WHERE user_id=? AND product_id=?", userID, productID)
	return n > 0, err
}

func (r *FavoriteRepository) CountByProduct(ctx context.Context, productID uint64) (int64, error) {
	return countCtx(ctx, r.conn, "SELECT COUNT(*) FROM product_favorites WHERE product_id=?", productID)
}

func (r *FavoriteRepository) List(ctx context.Context, userID uint64, page, pageSize int) ([]model.FavoriteListItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM product_favorites WHERE user_id=?", userID)
	if err != nil {
		return nil, 0, err
	}
	var favs []model.ProductFavorite
	if err := r.conn.QueryRowsPartialCtx(ctx, &favs,
		"SELECT "+productFavoriteColumns+" FROM product_favorites WHERE user_id=? ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?",
		userID, pageSize, (page-1)*pageSize,
	); err != nil {
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
	if err := r.conn.QueryRowsPartialCtx(ctx, &products,
		"SELECT id, name, IFNULL(main_image,'') AS main_image, sale_price, status, collect_count FROM products WHERE id IN ("+placeholders(len(ids))+")",
		inArgs(ids)...,
	); err != nil {
		return nil, 0, err
	}
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
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE products SET avg_rating=?, review_count=?, good_rate=? WHERE id=?",
		avg, count, goodRate, productID)
	return err
}
