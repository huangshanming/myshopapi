package repository

import (
	"context"
	"strings"

	"mymall/services/order-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const logisticsColumns = "id, name, code, sort, status, created_at, updated_at"

type LogisticsListFilter struct {
	Name        string
	Code        string
	Status      *int8
	Keyword     string
	Page        int
	PageSize    int
	EnabledOnly bool
}

type LogisticsRepository struct {
	conn sqlx.SqlConn
}

func NewLogisticsRepository(conn sqlx.SqlConn) *LogisticsRepository {
	return &LogisticsRepository{conn: conn}
}

func (r *LogisticsRepository) List(ctx context.Context, f LogisticsListFilter) ([]model.LogisticsCompany, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	where := []string{"1=1"}
	args := make([]any, 0, 8)
	if f.EnabledOnly {
		where = append(where, "status=1")
	} else if f.Status != nil {
		where = append(where, "status=?")
		args = append(args, *f.Status)
	}
	if f.Name != "" {
		where = append(where, "name LIKE ?")
		args = append(args, "%"+f.Name+"%")
	}
	if f.Code != "" {
		where = append(where, "code LIKE ?")
		args = append(args, "%"+f.Code+"%")
	}
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		where = append(where, "(name LIKE ? OR code LIKE ?)")
		args = append(args, like, like)
	}
	w := strings.Join(where, " AND ")
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM logistics_companies WHERE "+w, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), (f.Page-1)*f.PageSize, f.PageSize)
	var list []model.LogisticsCompany
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+logisticsColumns+" FROM logistics_companies WHERE "+w+" ORDER BY sort ASC, id ASC LIMIT ?, ?",
		listArgs...,
	)
	return list, total, err
}

func (r *LogisticsRepository) Options(ctx context.Context, keyword string, limit int) ([]model.LogisticsCompany, error) {
	if limit < 1 {
		limit = 50
	}
	where := []string{"status=1"}
	args := make([]any, 0, 3)
	if keyword != "" {
		like := "%" + keyword + "%"
		where = append(where, "(name LIKE ? OR code LIKE ?)")
		args = append(args, like, like)
	}
	args = append(args, limit)
	var list []model.LogisticsCompany
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+logisticsColumns+" FROM logistics_companies WHERE "+strings.Join(where, " AND ")+" ORDER BY sort ASC, id ASC LIMIT ?",
		args...,
	)
	return list, err
}

func (r *LogisticsRepository) Create(ctx context.Context, c *model.LogisticsCompany) error {
	id, err := lastInsertID(ctx, r.conn,
		"INSERT INTO logistics_companies (name, code, sort, status) VALUES (?, ?, ?, ?)",
		c.Name, c.Code, c.Sort, c.Status,
	)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *LogisticsRepository) Update(ctx context.Context, id uint64, name, code string, sort int) error {
	n, err := execAffected(ctx, r.conn,
		"UPDATE logistics_companies SET name=?, code=?, sort=? WHERE id=?",
		name, code, sort, id,
	)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *LogisticsRepository) UpdateStatus(ctx context.Context, id uint64, status int8) error {
	n, err := execAffected(ctx, r.conn, "UPDATE logistics_companies SET status=? WHERE id=?", status, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *LogisticsRepository) Delete(ctx context.Context, id uint64) error {
	n, err := execAffected(ctx, r.conn, "DELETE FROM logistics_companies WHERE id=?", id)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *LogisticsRepository) FindByCode(ctx context.Context, code string) (*model.LogisticsCompany, error) {
	var c model.LogisticsCompany
	err := r.conn.QueryRowCtx(ctx, &c, "SELECT "+logisticsColumns+" FROM logistics_companies WHERE code=? LIMIT 1", code)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *LogisticsRepository) SeedDefaults(ctx context.Context) error {
	seeds := []model.LogisticsCompany{
		{Name: "顺丰速运", Code: "SF", Sort: 10, Status: 1},
		{Name: "中通快递", Code: "ZTO", Sort: 20, Status: 1},
		{Name: "圆通速递", Code: "YTO", Sort: 30, Status: 1},
		{Name: "韵达快递", Code: "YD", Sort: 40, Status: 1},
		{Name: "申通快递", Code: "STO", Sort: 50, Status: 1},
		{Name: "EMS", Code: "EMS", Sort: 60, Status: 1},
		{Name: "京东物流", Code: "JD", Sort: 70, Status: 1},
		{Name: "德邦快递", Code: "DBL", Sort: 80, Status: 1},
	}
	for _, s := range seeds {
		n, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM logistics_companies WHERE code=?", s.Code)
		if err != nil {
			return err
		}
		if n == 0 {
			if err := r.Create(ctx, &s); err != nil {
				return err
			}
		}
	}
	return nil
}
