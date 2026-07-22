package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/user-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PointsOrderRepository struct {
	conn sqlx.SqlConn
}

func NewPointsOrderRepository(conn sqlx.SqlConn) *PointsOrderRepository {
	return &PointsOrderRepository{conn: conn}
}

type PointsOrderListFilter struct {
	Status  string
	OrderNo string
	UserID  uint64
	Keyword string
}

func (r *PointsOrderRepository) List(ctx context.Context, page, pageSize int, f PointsOrderListFilter) ([]model.PointsExchangeOrder, int64, error) {
	where := "1=1"
	args := make([]any, 0, 6)
	if f.Status != "" {
		where += " AND status=?"
		args = append(args, f.Status)
	}
	if f.OrderNo != "" {
		where += " AND order_no=?"
		args = append(args, strings.TrimSpace(f.OrderNo))
	}
	if f.UserID > 0 {
		where += " AND user_id=?"
		args = append(args, f.UserID)
	}
	if kw := strings.TrimSpace(f.Keyword); kw != "" {
		where += " AND (product_name LIKE ? OR order_no LIKE ?)"
		args = append(args, "%"+kw+"%", "%"+kw+"%")
	}
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM points_exchange_orders WHERE "+where, args...,
	)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.PointsExchangeOrder
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+pointsOrderColumns+" FROM points_exchange_orders WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *PointsOrderRepository) GetByID(ctx context.Context, id uint64) (*model.PointsExchangeOrder, error) {
	var o model.PointsExchangeOrder
	err := r.conn.QueryRowCtx(ctx, &o,
		"SELECT "+pointsOrderColumns+" FROM points_exchange_orders WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func genOrderNo() string {
	return fmt.Sprintf("PE%s%04d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%10000)
}

func (r *PointsOrderRepository) CreateExchangeLocal(ctx context.Context, userID, productID uint64, quantity int, receiverName, receiverPhone, receiverAddress string) (*model.PointsExchangeOrder, error) {
	if quantity < 1 {
		quantity = 1
	}
	var out *model.PointsExchangeOrder
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var p model.PointsProduct
		if err := session.QueryRowCtx(ctx, &p,
			"SELECT "+pointsProductColumns+" FROM points_products WHERE id=? FOR UPDATE", productID,
		); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return errors.New("商品不存在")
			}
			return err
		}
		if p.Status != model.PointsProductStatusOn {
			return errors.New("商品已下架")
		}
		if p.Stock < quantity {
			return errors.New("库存不足")
		}
		cost := p.PointsPrice * quantity
		if cost <= 0 {
			return errors.New("商品积分价无效")
		}
		if p.PerUserLimit > 0 {
			n, err := countQuery(ctx, session,
				"SELECT COUNT(*) FROM points_exchange_orders WHERE user_id=? AND product_id=? AND status<>?",
				userID, productID, model.PointsOrderCancelled,
			)
			if err != nil {
				return err
			}
			if int(n)+quantity > p.PerUserLimit {
				return errors.New("已达兑换上限")
			}
		}
		if _, err := session.ExecCtx(ctx,
			"UPDATE points_products SET stock=? WHERE id=?", p.Stock-quantity, p.ID,
		); err != nil {
			return err
		}
		res, err := session.ExecCtx(ctx,
			`INSERT INTO points_exchange_orders (order_no, user_id, product_id, product_name, product_cover, quantity, points_cost, status, receiver_name, receiver_phone, receiver_address)
			 VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
			genOrderNo(), userID, p.ID, p.Name, p.CoverURL, quantity, cost, model.PointsOrderPending,
			strings.TrimSpace(receiverName), strings.TrimSpace(receiverPhone), strings.TrimSpace(receiverAddress),
		)
		if err != nil {
			return err
		}
		id, err := lastInsertID(res)
		if err != nil {
			return err
		}
		o := &model.PointsExchangeOrder{}
		if err := session.QueryRowCtx(ctx, o,
			"SELECT "+pointsOrderColumns+" FROM points_exchange_orders WHERE id=? LIMIT 1", id,
		); err != nil {
			return err
		}
		out = o
		return nil
	})
	return out, err
}

func (r *PointsOrderRepository) AbortExchange(ctx context.Context, id uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var o model.PointsExchangeOrder
		if err := session.QueryRowCtx(ctx, &o,
			"SELECT "+pointsOrderColumns+" FROM points_exchange_orders WHERE id=? FOR UPDATE", id,
		); err != nil {
			return err
		}
		if o.Status != model.PointsOrderPending {
			return nil
		}
		if _, err := session.ExecCtx(ctx,
			"UPDATE points_products SET stock=stock+? WHERE id=?", o.Quantity, o.ProductID,
		); err != nil {
			return err
		}
		_, err := session.ExecCtx(ctx, "DELETE FROM points_exchange_orders WHERE id=?", id)
		return err
	})
}

func (r *PointsOrderRepository) Ship(ctx context.Context, id uint64, company, shipNo string) (*model.PointsExchangeOrder, error) {
	o, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if o.Status != model.PointsOrderPending {
		return nil, errors.New("当前状态不可发货")
	}
	now := common.LocalTime(time.Now())
	if _, err := r.conn.ExecCtx(ctx,
		"UPDATE points_exchange_orders SET status=?, ship_company=?, ship_no=?, shipped_at=? WHERE id=?",
		model.PointsOrderShipped, strings.TrimSpace(company), strings.TrimSpace(shipNo), &now, id,
	); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *PointsOrderRepository) Complete(ctx context.Context, id uint64) (*model.PointsExchangeOrder, error) {
	o, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if o.Status != model.PointsOrderPending && o.Status != model.PointsOrderShipped {
		return nil, errors.New("当前状态不可完成")
	}
	now := common.LocalTime(time.Now())
	if _, err := r.conn.ExecCtx(ctx,
		"UPDATE points_exchange_orders SET status=?, completed_at=? WHERE id=?",
		model.PointsOrderCompleted, &now, id,
	); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *PointsOrderRepository) CancelLocal(ctx context.Context, id uint64, remark string) (*model.PointsExchangeOrder, error) {
	var out *model.PointsExchangeOrder
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var o model.PointsExchangeOrder
		if err := session.QueryRowCtx(ctx, &o,
			"SELECT "+pointsOrderColumns+" FROM points_exchange_orders WHERE id=? FOR UPDATE", id,
		); err != nil {
			return errors.New("订单不存在")
		}
		if o.Status != model.PointsOrderPending {
			return errors.New("仅待发货订单可取消退积分")
		}
		now := common.LocalTime(time.Now())
		adminRemark := o.AdminRemark
		if strings.TrimSpace(remark) != "" {
			adminRemark = strings.TrimSpace(remark)
		}
		if _, err := session.ExecCtx(ctx,
			"UPDATE points_exchange_orders SET status=?, cancelled_at=?, admin_remark=? WHERE id=?",
			model.PointsOrderCancelled, &now, adminRemark, id,
		); err != nil {
			return err
		}
		if _, err := session.ExecCtx(ctx,
			"UPDATE points_products SET stock=stock+? WHERE id=?", o.Quantity, o.ProductID,
		); err != nil {
			return err
		}
		out = &o
		out.Status = model.PointsOrderCancelled
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *PointsOrderRepository) Remark(ctx context.Context, id uint64, remark string) (*model.PointsExchangeOrder, error) {
	if _, err := r.GetByID(ctx, id); err != nil {
		return nil, errors.New("订单不存在")
	}
	if _, err := r.conn.ExecCtx(ctx,
		"UPDATE points_exchange_orders SET admin_remark=? WHERE id=?",
		strings.TrimSpace(remark), id,
	); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *PointsOrderRepository) MapUserBriefs(ctx context.Context, ids []uint64) map[uint64][2]string {
	out := map[uint64][2]string{}
	if len(ids) == 0 {
		return out
	}
	type row struct {
		ID       uint64 `db:"id"`
		Nickname string `db:"nickname"`
		Mobile   string `db:"mobile"`
	}
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(
		"SELECT id, nickname, mobile FROM users WHERE id IN (%s) AND deleted_at IS NULL",
		strings.Join(placeholders, ","),
	)
	var rows []row
	if err := r.conn.QueryRowsCtx(ctx, &rows, query, args...); err != nil {
		return out
	}
	for _, u := range rows {
		out[u.ID] = [2]string{u.Nickname, u.Mobile}
	}
	return out
}
