package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"mymall/services/lottery-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LotteryRepository struct {
	conn sqlx.SqlConn
}

func NewLotteryRepository(conn sqlx.SqlConn) *LotteryRepository {
	return &LotteryRepository{conn: conn}
}

func (r *LotteryRepository) Conn() sqlx.SqlConn { return r.conn }

const activityCols = "id, title, status, cost_points, daily_limit, start_at, end_at, created_at, updated_at"
const prizeCols = "id, activity_id, slot, name, cover_url, prize_type, points_amount, weight, stock, stock_strict, created_at, updated_at"
const recordCols = "id, user_id, activity_id, prize_id, slot, prize_type, prize_name, points_amount, cost_points, status, fulfill_status, address_id, receiver_name, receiver_phone, receiver_address, ship_company, ship_no, shipped_at, created_at"

func (r *LotteryRepository) GetActiveActivity(ctx context.Context) (*model.LotteryActivity, error) {
	var a model.LotteryActivity
	now := time.Now()
	err := r.conn.QueryRowPartialCtx(ctx, &a,
		`SELECT `+activityCols+` FROM lottery_activities
		 WHERE status=? AND (start_at IS NULL OR start_at<=?) AND (end_at IS NULL OR end_at>=?)
		 ORDER BY id DESC LIMIT 1`,
		model.ActivityStatusOnline, now, now,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *LotteryRepository) GetActivity(ctx context.Context, id uint64) (*model.LotteryActivity, error) {
	var a model.LotteryActivity
	err := r.conn.QueryRowPartialCtx(ctx, &a,
		`SELECT `+activityCols+` FROM lottery_activities WHERE id=? LIMIT 1`, id,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *LotteryRepository) ListActivities(ctx context.Context, page, pageSize int) ([]model.LotteryActivity, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var total int64
	if err := r.conn.QueryRowCtx(ctx, &total, `SELECT COUNT(*) FROM lottery_activities`); err != nil {
		return nil, 0, err
	}
	var list []model.LotteryActivity
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		`SELECT `+activityCols+` FROM lottery_activities ORDER BY id DESC LIMIT ? OFFSET ?`,
		pageSize, (page-1)*pageSize,
	)
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		list = []model.LotteryActivity{}
	}
	return list, total, nil
}

func (r *LotteryRepository) CreateActivity(ctx context.Context, a *model.LotteryActivity) (uint64, error) {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO lottery_activities (title, status, cost_points, daily_limit, start_at, end_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		a.Title, a.Status, a.CostPoints, a.DailyLimit, a.StartAt, a.EndAt,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *LotteryRepository) UpdateActivity(ctx context.Context, a *model.LotteryActivity) error {
	_, err := r.conn.ExecCtx(ctx,
		`UPDATE lottery_activities SET title=?, status=?, cost_points=?, daily_limit=?, start_at=?, end_at=? WHERE id=?`,
		a.Title, a.Status, a.CostPoints, a.DailyLimit, a.StartAt, a.EndAt, a.ID,
	)
	return err
}

func (r *LotteryRepository) ListPrizes(ctx context.Context, activityID uint64) ([]model.LotteryPrize, error) {
	var list []model.LotteryPrize
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		`SELECT `+prizeCols+` FROM lottery_prizes WHERE activity_id=? ORDER BY slot ASC`, activityID,
	)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []model.LotteryPrize{}
	}
	return list, nil
}

func (r *LotteryRepository) ReplacePrizes(ctx context.Context, activityID uint64, prizes []model.LotteryPrize) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, `DELETE FROM lottery_prizes WHERE activity_id=?`, activityID); err != nil {
			return err
		}
		for _, p := range prizes {
			if _, err := session.ExecCtx(ctx,
				`INSERT INTO lottery_prizes (activity_id, slot, name, cover_url, prize_type, points_amount, weight, stock, stock_strict)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				activityID, p.Slot, p.Name, p.CoverURL, p.PrizeType, p.PointsAmount, p.Weight, p.Stock, p.StockStrict,
			); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *LotteryRepository) CountUserDrawsToday(ctx context.Context, userID, activityID uint64) (int64, error) {
	start := time.Now().Truncate(24 * time.Hour)
	var n int64
	err := r.conn.QueryRowCtx(ctx, &n,
		`SELECT COUNT(*) FROM lottery_draw_records
		 WHERE user_id=? AND activity_id=? AND status=? AND created_at>=?`,
		userID, activityID, model.RecordStatusDone, start,
	)
	return n, err
}

func (r *LotteryRepository) InsertPendingRecord(ctx context.Context, rec *model.LotteryDrawRecord) (uint64, error) {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO lottery_draw_records
		 (user_id, activity_id, prize_id, slot, prize_type, prize_name, points_amount, cost_points, status, fulfill_status)
		 VALUES (?, ?, 0, 0, '', '', 0, ?, ?, ?)`,
		rec.UserID, rec.ActivityID, rec.CostPoints, model.RecordStatusPending, model.FulfillNone,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *LotteryRepository) MarkRecordFailed(ctx context.Context, id uint64) error {
	_, err := r.conn.ExecCtx(ctx,
		`UPDATE lottery_draw_records SET status=? WHERE id=?`, model.RecordStatusFailed, id,
	)
	return err
}

// FinalizeDraw picks prize by weight, decrements stock, updates record — all in one TX.
func (r *LotteryRepository) FinalizeDraw(ctx context.Context, recordID uint64, prizes []model.LotteryPrize, pickIdx int) (*model.LotteryPrize, error) {
	if pickIdx < 0 || pickIdx >= len(prizes) {
		return nil, errors.New("抽奖结果无效")
	}
	chosen := prizes[pickIdx]
	var out model.LotteryPrize
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// re-lock prize row
		var p model.LotteryPrize
		err := session.QueryRowPartialCtx(ctx, &p,
			`SELECT `+prizeCols+` FROM lottery_prizes WHERE id=? FOR UPDATE`, chosen.ID,
		)
		if err != nil {
			return err
		}
		if p.Weight <= 0 {
			return errors.New("奖品已不可用")
		}
		if p.Stock == 0 {
			return errors.New("奖品已抽完")
		}
		if p.Stock > 0 {
			res, err := session.ExecCtx(ctx,
				`UPDATE lottery_prizes SET stock=stock-1 WHERE id=? AND stock>0`, p.ID,
			)
			if err != nil {
				return err
			}
			n, _ := res.RowsAffected()
			if n == 0 {
				return errors.New("奖品已抽完")
			}
			p.Stock--
		}
		fulfill := model.FulfillNone
		if p.PrizeType == model.PrizeTypePhysical {
			fulfill = model.FulfillNeedAddress
		}
		_, err = session.ExecCtx(ctx,
			`UPDATE lottery_draw_records SET prize_id=?, slot=?, prize_type=?, prize_name=?, points_amount=?, status=?, fulfill_status=? WHERE id=? AND status=?`,
			p.ID, p.Slot, p.PrizeType, p.Name, p.PointsAmount, model.RecordStatusDone, fulfill, recordID, model.RecordStatusPending,
		)
		if err != nil {
			return err
		}
		out = p
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *LotteryRepository) GetRecord(ctx context.Context, id uint64) (*model.LotteryDrawRecord, error) {
	var rec model.LotteryDrawRecord
	err := r.conn.QueryRowPartialCtx(ctx, &rec,
		`SELECT `+recordCols+` FROM lottery_draw_records WHERE id=? LIMIT 1`, id,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *LotteryRepository) ClaimAddress(ctx context.Context, id, userID, addressID uint64, name, phone, address string) error {
	res, err := r.conn.ExecCtx(ctx,
		`UPDATE lottery_draw_records
		 SET address_id=?, receiver_name=?, receiver_phone=?, receiver_address=?, fulfill_status=?
		 WHERE id=? AND user_id=? AND status=? AND fulfill_status=? AND prize_type=?`,
		addressID, name, phone, address, model.FulfillPending,
		id, userID, model.RecordStatusDone, model.FulfillNeedAddress, model.PrizeTypePhysical,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("记录不存在或无需填写地址")
	}
	return nil
}

func (r *LotteryRepository) ShipRecord(ctx context.Context, id uint64, company, shipNo string) error {
	now := time.Now()
	res, err := r.conn.ExecCtx(ctx,
		`UPDATE lottery_draw_records
		 SET ship_company=?, ship_no=?, shipped_at=?, fulfill_status=?
		 WHERE id=? AND status=? AND fulfill_status=? AND prize_type=?`,
		company, shipNo, now, model.FulfillShipped,
		id, model.RecordStatusDone, model.FulfillPending, model.PrizeTypePhysical,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("订单不存在或状态不可发货")
	}
	return nil
}

func (r *LotteryRepository) ListFulfillmentOrders(ctx context.Context, fulfillStatus string, page, pageSize int) ([]model.LotteryDrawRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	where := []string{"status=?", "prize_type=?"}
	args := []any{model.RecordStatusDone, model.PrizeTypePhysical}
	if s := strings.TrimSpace(fulfillStatus); s != "" {
		where = append(where, "fulfill_status=?")
		args = append(args, s)
	} else {
		where = append(where, "fulfill_status IN (?,?,?)")
		args = append(args, model.FulfillNeedAddress, model.FulfillPending, model.FulfillShipped)
	}
	w := strings.Join(where, " AND ")
	var total int64
	if err := r.conn.QueryRowCtx(ctx, &total, `SELECT COUNT(*) FROM lottery_draw_records WHERE `+w, args...); err != nil {
		return nil, 0, err
	}
	args2 := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.LotteryDrawRecord
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		`SELECT `+recordCols+` FROM lottery_draw_records WHERE `+w+` ORDER BY id DESC LIMIT ? OFFSET ?`,
		args2...,
	)
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		list = []model.LotteryDrawRecord{}
	}
	return list, total, nil
}

func (r *LotteryRepository) ListUserRecords(ctx context.Context, userID uint64, page, pageSize int) ([]model.LotteryDrawRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var total int64
	if err := r.conn.QueryRowCtx(ctx, &total,
		`SELECT COUNT(*) FROM lottery_draw_records WHERE user_id=? AND status=?`,
		userID, model.RecordStatusDone,
	); err != nil {
		return nil, 0, err
	}
	var list []model.LotteryDrawRecord
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		`SELECT `+recordCols+` FROM lottery_draw_records WHERE user_id=? AND status=? ORDER BY id DESC LIMIT ? OFFSET ?`,
		userID, model.RecordStatusDone, pageSize, (page-1)*pageSize,
	)
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		list = []model.LotteryDrawRecord{}
	}
	return list, total, nil
}

func (r *LotteryRepository) ListRecordsAdmin(ctx context.Context, activityID uint64, prizeType string, page, pageSize int) ([]model.LotteryDrawRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	where := []string{"status=?"}
	args := []any{model.RecordStatusDone}
	if activityID > 0 {
		where = append(where, "activity_id=?")
		args = append(args, activityID)
	}
	if pt := strings.TrimSpace(prizeType); pt != "" {
		where = append(where, "prize_type=?")
		args = append(args, pt)
	}
	w := strings.Join(where, " AND ")
	var total int64
	if err := r.conn.QueryRowCtx(ctx, &total, `SELECT COUNT(*) FROM lottery_draw_records WHERE `+w, args...); err != nil {
		return nil, 0, err
	}
	args2 := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.LotteryDrawRecord
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		`SELECT `+recordCols+` FROM lottery_draw_records WHERE `+w+` ORDER BY id DESC LIMIT ? OFFSET ?`,
		args2...,
	)
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		list = []model.LotteryDrawRecord{}
	}
	return list, total, nil
}

func ParseTimePtr(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			tt := t
			return &tt, nil
		}
	}
	return nil, fmt.Errorf("时间格式无效: %s", s)
}

func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// EnsureNullTime helps sql.NullTime unused import if needed.
var _ = sql.ErrNoRows
