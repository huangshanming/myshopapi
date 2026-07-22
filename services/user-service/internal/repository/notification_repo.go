package repository

import (
	"context"
	"errors"

	"mymall/services/user-service/internal/model"
)

const notificationColumns = "id, user_id, IFNULL(title,'') AS title, IFNULL(content,'') AS content, IFNULL(msg_type,'') AS msg_type, IFNULL(link_type,'') AS link_type, IFNULL(link_id,0) AS link_id, IFNULL(extra,'') AS extra, is_read, IFNULL(sender_type,'') AS sender_type, IFNULL(sender_id,0) AS sender_id, IFNULL(batch_id,0) AS batch_id, created_at"
const notificationBatchColumns = "id, IFNULL(title,'') AS title, IFNULL(content,'') AS content, IFNULL(target,'') AS target, user_count, success_count, IFNULL(link_type,'') AS link_type, IFNULL(link_id,0) AS link_id, IFNULL(sender_id,0) AS sender_id, created_at"

func (r *UserRepository) CreateNotification(ctx context.Context, n *model.UserNotification) error {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO user_notifications (user_id, title, content, msg_type, link_type, link_id, extra, is_read, sender_type, sender_id, batch_id)
		 VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		n.UserID, n.Title, n.Content, n.MsgType, n.LinkType, n.LinkID, n.Extra, n.IsRead, n.SenderType, n.SenderID, n.BatchID,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	n.ID = id
	return nil
}

func (r *UserRepository) CreateNotifications(ctx context.Context, list []model.UserNotification) error {
	if len(list) == 0 {
		return nil
	}
	const batchSize = 200
	for i := 0; i < len(list); i += batchSize {
		end := i + batchSize
		if end > len(list) {
			end = len(list)
		}
		chunk := list[i:end]
		for j := range chunk {
			if err := r.CreateNotification(ctx, &chunk[j]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *UserRepository) ListNotifications(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserNotification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_notifications WHERE user_id=?", userID,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.UserNotification
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+notificationColumns+" FROM user_notifications WHERE user_id=? ORDER BY id DESC LIMIT ? OFFSET ?",
		userID, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *UserRepository) UnreadNotificationCount(ctx context.Context, userID uint64) (int64, error) {
	return countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_notifications WHERE user_id=? AND is_read=0", userID,
	)
}

func (r *UserRepository) MarkNotificationRead(ctx context.Context, userID, id uint64) error {
	n, err := execRows(ctx, r.conn,
		"UPDATE user_notifications SET is_read=1 WHERE id=? AND user_id=?", id, userID,
	)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("消息不存在")
	}
	return nil
}

func (r *UserRepository) MarkAllNotificationsRead(ctx context.Context, userID uint64) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE user_notifications SET is_read=1 WHERE user_id=? AND is_read=0", userID,
	)
	return err
}

func (r *UserRepository) ListActiveUserIDs(ctx context.Context, offset, limit int) ([]uint64, error) {
	type row struct {
		ID uint64 `db:"id"`
	}
	var rows []row
	err := r.conn.QueryRowsPartialCtx(ctx, &rows,
		"SELECT id FROM users WHERE status=1 AND deleted_at IS NULL ORDER BY id ASC LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(rows))
	for _, item := range rows {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

func (r *UserRepository) CreateNotificationBatch(ctx context.Context, b *model.UserNotificationBatch) error {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO user_notification_batches (title, content, target, user_count, success_count, link_type, link_id, sender_id)
		 VALUES (?,?,?,?,?,?,?,?)`,
		b.Title, b.Content, b.Target, b.UserCount, b.SuccessCount, b.LinkType, b.LinkID, b.SenderID,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}

func (r *UserRepository) UpdateNotificationBatchSuccess(ctx context.Context, id uint64, success int) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE user_notification_batches SET success_count=? WHERE id=?", success, id,
	)
	return err
}

func (r *UserRepository) ListNotificationBatches(ctx context.Context, page, pageSize int) ([]model.UserNotificationBatch, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM user_notification_batches")
	if err != nil {
		return nil, 0, err
	}
	var list []model.UserNotificationBatch
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+notificationBatchColumns+" FROM user_notification_batches ORDER BY id DESC LIMIT ? OFFSET ?",
		pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *UserRepository) GetNotificationBatch(ctx context.Context, id uint64) (*model.UserNotificationBatch, error) {
	var b model.UserNotificationBatch
	err := r.conn.QueryRowPartialCtx(ctx, &b,
		"SELECT "+notificationBatchColumns+" FROM user_notification_batches WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

type BatchRecipientRow struct {
	UserID   uint64 `db:"user_id" json:"user_id"`
	Nickname string `db:"nickname" json:"nickname"`
	Mobile   string `db:"mobile" json:"mobile"`
	IsRead   int8   `db:"is_read" json:"is_read"`
}

func (r *UserRepository) ListBatchRecipients(ctx context.Context, batchID uint64, page, pageSize int) ([]BatchRecipientRow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_notifications WHERE batch_id=?", batchID,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []BatchRecipientRow
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		`SELECT n.user_id, COALESCE(u.nickname,'') AS nickname, COALESCE(u.mobile,'') AS mobile, n.is_read
		 FROM user_notifications n
		 LEFT JOIN users u ON u.id = n.user_id
		 WHERE n.batch_id=? ORDER BY n.id ASC LIMIT ? OFFSET ?`,
		batchID, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}
