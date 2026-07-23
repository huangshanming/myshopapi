package repository

import (
	"context"
	"strings"

	"mymall/services/user-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const pointsProductColumns = "id, IFNULL(name,'') AS name, IFNULL(cover_url,'') AS cover_url, IFNULL(description,'') AS description, points_price, stock, per_user_limit, status, sort, created_at, updated_at"

type PointsProductRepository struct {
	conn sqlx.SqlConn
}

func NewPointsProductRepository(conn sqlx.SqlConn) *PointsProductRepository {
	return &PointsProductRepository{conn: conn}
}

func (r *PointsProductRepository) List(ctx context.Context, page, pageSize int, status, keyword string) ([]model.PointsProduct, int64, error) {
	where := "1=1"
	args := make([]any, 0, 4)
	if status != "" {
		where += " AND status=?"
		args = append(args, status)
	}
	if kw := strings.TrimSpace(keyword); kw != "" {
		where += " AND name LIKE ?"
		args = append(args, "%"+kw+"%")
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM points_products WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	list := make([]model.PointsProduct, 0)
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+pointsProductColumns+" FROM points_products WHERE "+where+" ORDER BY sort DESC, id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *PointsProductRepository) GetByID(ctx context.Context, id uint64) (*model.PointsProduct, error) {
	var p model.PointsProduct
	err := r.conn.QueryRowPartialCtx(ctx, &p,
		"SELECT "+pointsProductColumns+" FROM points_products WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PointsProductRepository) Create(ctx context.Context, p *model.PointsProduct) error {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO points_products (name, cover_url, description, points_price, stock, per_user_limit, status, sort)
		 VALUES (?,?,?,?,?,?,?,?)`,
		p.Name, p.CoverURL, p.Description, p.PointsPrice, p.Stock, p.PerUserLimit, p.Status, p.Sort,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func (r *PointsProductRepository) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("points_products", updates, "id=?", id)
	if err != nil {
		return err
	}
	_, err = r.conn.ExecCtx(ctx, query, args...)
	return err
}

func (r *PointsProductRepository) Delete(ctx context.Context, id uint64) error {
	_, err := r.conn.ExecCtx(ctx, "DELETE FROM points_products WHERE id=?", id)
	return err
}
