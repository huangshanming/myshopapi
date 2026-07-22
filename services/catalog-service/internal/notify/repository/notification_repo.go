package repository

import (
	"context"
	"errors"
	"strings"

	"mymall/services/catalog-service/internal/notify/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const shopNotificationColumns = "id, shop_id, type, title, content, link, ref_type, ref_id, is_read, created_at"

type NotificationRepository struct {
	conn sqlx.SqlConn
}

func NewNotificationRepository(conn sqlx.SqlConn) *NotificationRepository {
	return &NotificationRepository{conn: conn}
}

func (r *NotificationRepository) Create(ctx context.Context, n *model.ShopNotification) error {
	id, err := lastInsertID(ctx, r.conn,
		`INSERT INTO shop_notifications (shop_id, type, title, content, link, ref_type, ref_id, is_read)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		n.ShopID, n.Type, n.Title, n.Content, n.Link, n.RefType, n.RefID, n.IsRead,
	)
	if err != nil {
		return err
	}
	n.ID = id
	return nil
}

type NotificationListFilter struct {
	ShopID   uint64
	IsRead   *int8 // nil=全部
	Page     int
	PageSize int
}

func (r *NotificationRepository) List(ctx context.Context, f NotificationListFilter) ([]model.ShopNotification, int64, error) {
	where := []string{"shop_id=?"}
	args := []any{f.ShopID}
	if f.IsRead != nil {
		where = append(where, "is_read=?")
		args = append(args, *f.IsRead)
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM shop_notifications WHERE "+whereSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	var list []model.ShopNotification
	qArgs := append(args, f.PageSize, (f.Page-1)*f.PageSize)
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+shopNotificationColumns+" FROM shop_notifications WHERE "+whereSQL+" ORDER BY id DESC LIMIT ? OFFSET ?",
		qArgs...,
	)
	return list, total, err
}

func (r *NotificationRepository) UnreadCount(ctx context.Context, shopID uint64) (int64, error) {
	return countCtx(ctx, r.conn,
		"SELECT COUNT(*) FROM shop_notifications WHERE shop_id=? AND is_read=0", shopID)
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id, shopID uint64) error {
	n, err := execAffected(ctx, r.conn,
		"UPDATE shop_notifications SET is_read=1 WHERE id=? AND shop_id=?", id, shopID)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("通知不存在")
	}
	return nil
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, shopID uint64) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE shop_notifications SET is_read=1 WHERE shop_id=? AND is_read=0", shopID)
	return err
}
