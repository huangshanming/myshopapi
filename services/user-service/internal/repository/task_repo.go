package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mymall/common"
	"mymall/services/user-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	taskDefinitionColumns = "id, code, title, description, icon, period, enabled, reward_points, target_count, daily_limit, sort, rules_json, created_at, updated_at"
	taskProgressColumns   = "id, user_id, task_code, biz_date, progress, claim_count, status, claimed_at, created_at, updated_at"
	userPointsColumns     = "user_id, points, created_at, updated_at"
	userPointLogColumns   = "id, user_id, change_type, delta, points_after, remark, ref_type, ref_id, created_at"
	pointsOrderColumns    = "id, order_no, user_id, product_id, product_name, product_cover, quantity, points_cost, status, receiver_name, receiver_phone, receiver_address, ship_company, ship_no, admin_remark, shipped_at, completed_at, cancelled_at, created_at, updated_at"
)

type TaskRepository struct {
	conn sqlx.SqlConn
}

func NewTaskRepository(conn sqlx.SqlConn) *TaskRepository {
	return &TaskRepository{conn: conn}
}

func (r *TaskRepository) SeedIfEmpty(ctx context.Context) error {
	n, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM task_definitions")
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	seeds := []model.TaskDefinition{
		{Code: "daily_checkin", Title: "每日签到", Description: "每天签到一次可领取积分", Icon: "checkin", Period: model.TaskPeriodDaily, Enabled: 1, RewardPoints: 5, TargetCount: 1, DailyLimit: 1, Sort: 10, RulesJSON: "{}"},
		{Code: "publish_article", Title: "发布种草笔记", Description: "发布笔记并通过审核后可领取；可配置每日奖励次数", Icon: "article", Period: model.TaskPeriodDaily, Enabled: 1, RewardPoints: 20, TargetCount: 1, DailyLimit: 1, Sort: 20, RulesJSON: "{}"},
		{Code: "comment_article", Title: "发表评论", Description: "在种草社区发表有效评论", Icon: "comment", Period: model.TaskPeriodDaily, Enabled: 1, RewardPoints: 5, TargetCount: 3, DailyLimit: 1, Sort: 30, RulesJSON: "{}"},
		{Code: "like_article", Title: "点赞文章", Description: "为喜欢的笔记点赞", Icon: "like", Period: model.TaskPeriodDaily, Enabled: 1, RewardPoints: 3, TargetCount: 5, DailyLimit: 1, Sort: 40, RulesJSON: "{}"},
		{Code: "favorite_article", Title: "收藏文章", Description: "收藏种草笔记", Icon: "star", Period: model.TaskPeriodDaily, Enabled: 1, RewardPoints: 3, TargetCount: 3, DailyLimit: 1, Sort: 50, RulesJSON: "{}"},
		{Code: "browse_products", Title: "浏览商品", Description: "浏览不同商品详情页", Icon: "browse", Period: model.TaskPeriodDaily, Enabled: 1, RewardPoints: 5, TargetCount: 5, DailyLimit: 1, Sort: 60, RulesJSON: "{}"},
		{Code: "place_order", Title: "下单购物", Description: "支付成功一笔订单", Icon: "order", Period: model.TaskPeriodDaily, Enabled: 1, RewardPoints: 30, TargetCount: 1, DailyLimit: 1, Sort: 70, RulesJSON: "{}"},
		{Code: "first_profile", Title: "完善资料", Description: "完善昵称与头像", Icon: "profile", Period: model.TaskPeriodOnce, Enabled: 1, RewardPoints: 50, TargetCount: 1, DailyLimit: 0, Sort: 80, RulesJSON: "{}"},
		{Code: "first_favorite_product", Title: "首次收藏商品", Description: "完成第一次商品收藏", Icon: "favorite", Period: model.TaskPeriodOnce, Enabled: 1, RewardPoints: 20, TargetCount: 1, DailyLimit: 0, Sort: 90, RulesJSON: "{}"},
		{Code: "invite_placeholder", Title: "邀请好友", Description: "邀请好友注册（即将开放）", Icon: "invite", Period: model.TaskPeriodDaily, Enabled: 0, RewardPoints: 100, TargetCount: 1, DailyLimit: 1, Sort: 100, RulesJSON: "{}"},
	}
	for i := range seeds {
		s := &seeds[i]
		if _, err := r.conn.ExecCtx(ctx,
			`INSERT INTO task_definitions (code, title, description, icon, period, enabled, reward_points, target_count, daily_limit, sort, rules_json)
			 VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
			s.Code, s.Title, s.Description, s.Icon, s.Period, s.Enabled, s.RewardPoints, s.TargetCount, s.DailyLimit, s.Sort, s.RulesJSON,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *TaskRepository) ListDefinitions(ctx context.Context, all bool) ([]model.TaskDefinition, error) {
	var list []model.TaskDefinition
	if all {
		err := r.conn.QueryRowsCtx(ctx, &list,
			"SELECT "+taskDefinitionColumns+" FROM task_definitions ORDER BY sort ASC, id ASC",
		)
		return list, err
	}
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+taskDefinitionColumns+" FROM task_definitions WHERE enabled=1 ORDER BY sort ASC, id ASC",
	)
	return list, err
}

func (r *TaskRepository) GetDefinition(ctx context.Context, code string) (*model.TaskDefinition, error) {
	var t model.TaskDefinition
	err := r.conn.QueryRowCtx(ctx, &t,
		"SELECT "+taskDefinitionColumns+" FROM task_definitions WHERE code=? LIMIT 1", code,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) GetDefinitionByID(ctx context.Context, id uint64) (*model.TaskDefinition, error) {
	var t model.TaskDefinition
	err := r.conn.QueryRowCtx(ctx, &t,
		"SELECT "+taskDefinitionColumns+" FROM task_definitions WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) UpdateDefinition(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("task_definitions", updates, "id=?", id)
	if err != nil {
		return err
	}
	_, err = r.conn.ExecCtx(ctx, query, args...)
	return err
}

func TodayBizDate() string {
	return time.Now().Format("2006-01-02")
}

func BizDateFor(def *model.TaskDefinition) string {
	if def.Period == model.TaskPeriodOnce {
		return model.TaskBizDateOnce
	}
	return TodayBizDate()
}

func (r *TaskRepository) GetOrCreateProgress(ctx context.Context, userID uint64, def *model.TaskDefinition) (*model.UserTaskProgress, error) {
	biz := BizDateFor(def)
	var p model.UserTaskProgress
	err := r.conn.QueryRowCtx(ctx, &p,
		"SELECT "+taskProgressColumns+" FROM user_task_progress WHERE user_id=? AND task_code=? AND biz_date=? LIMIT 1",
		userID, def.Code, biz,
	)
	if err == nil {
		return &p, nil
	}
	if !errors.Is(err, sqlx.ErrNotFound) {
		return nil, err
	}
	_, err = r.conn.ExecCtx(ctx,
		`INSERT IGNORE INTO user_task_progress (user_id, task_code, biz_date, progress, claim_count, status)
		 VALUES (?,?,?,0,0,?)`,
		userID, def.Code, biz, model.TaskStatusOngoing,
	)
	if err != nil {
		return nil, err
	}
	err = r.conn.QueryRowCtx(ctx, &p,
		"SELECT "+taskProgressColumns+" FROM user_task_progress WHERE user_id=? AND task_code=? AND biz_date=? LIMIT 1",
		userID, def.Code, biz,
	)
	return &p, err
}

func (r *TaskRepository) TryDedupe(ctx context.Context, userID uint64, taskCode, bizDate, refKey string) (bool, error) {
	if refKey == "" {
		return true, nil
	}
	res, err := r.conn.ExecCtx(ctx,
		"INSERT IGNORE INTO user_task_dedupe (user_id, task_code, biz_date, ref_key) VALUES (?,?,?,?)",
		userID, taskCode, bizDate, refKey,
	)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

func (r *TaskRepository) ApplyEvent(ctx context.Context, userID uint64, def *model.TaskDefinition, delta int, refKey string) (*model.UserTaskProgress, error) {
	if def.Enabled != 1 {
		return nil, nil
	}
	if delta < 1 {
		delta = 1
	}
	biz := BizDateFor(def)
	ok, err := r.TryDedupe(ctx, userID, def.Code, biz, refKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		return r.GetOrCreateProgress(ctx, userID, def)
	}

	var out *model.UserTaskProgress
	err = r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var p model.UserTaskProgress
		e := session.QueryRowCtx(ctx, &p,
			"SELECT "+taskProgressColumns+" FROM user_task_progress WHERE user_id=? AND task_code=? AND biz_date=? FOR UPDATE",
			userID, def.Code, biz,
		)
		if errors.Is(e, sqlx.ErrNotFound) {
			res, err := session.ExecCtx(ctx,
				`INSERT INTO user_task_progress (user_id, task_code, biz_date, progress, claim_count, status)
				 VALUES (?,?,?,0,0,?)`,
				userID, def.Code, biz, model.TaskStatusOngoing,
			)
			if err != nil {
				return err
			}
			id, err := lastInsertID(res)
			if err != nil {
				return err
			}
			e = session.QueryRowCtx(ctx, &p,
				"SELECT "+taskProgressColumns+" FROM user_task_progress WHERE id=? FOR UPDATE", id,
			)
		}
		if e != nil {
			return e
		}
		if def.Period == model.TaskPeriodOnce && p.Status == model.TaskStatusClaimed {
			out = &p
			return nil
		}
		if def.Period == model.TaskPeriodDaily && def.DailyLimit > 0 && p.ClaimCount >= def.DailyLimit && p.Status != model.TaskStatusClaimable {
			out = &p
			return nil
		}
		if p.Status == model.TaskStatusClaimable {
			out = &p
			return nil
		}
		if p.Status == model.TaskStatusClaimed && def.Period == model.TaskPeriodDaily {
			if def.DailyLimit == 0 || p.ClaimCount < def.DailyLimit {
				p.Progress = 0
				p.Status = model.TaskStatusOngoing
			} else {
				out = &p
				return nil
			}
		}
		target := def.TargetCount
		if target < 1 {
			target = 1
		}
		p.Progress += delta
		if p.Progress >= target {
			p.Progress = target
			p.Status = model.TaskStatusClaimable
		}
		if _, err := session.ExecCtx(ctx,
			"UPDATE user_task_progress SET progress=?, status=? WHERE id=?",
			p.Progress, p.Status, p.ID,
		); err != nil {
			return err
		}
		out = &p
		return nil
	})
	return out, err
}

func (r *TaskRepository) Claim(ctx context.Context, userID uint64, def *model.TaskDefinition) (*model.UserPoints, error) {
	biz := BizDateFor(def)
	var points *model.UserPoints
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var p model.UserTaskProgress
		if err := session.QueryRowCtx(ctx, &p,
			"SELECT "+taskProgressColumns+" FROM user_task_progress WHERE user_id=? AND task_code=? AND biz_date=? FOR UPDATE",
			userID, def.Code, biz,
		); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return errors.New("任务尚未完成")
			}
			return err
		}
		if p.Status != model.TaskStatusClaimable {
			if p.Status == model.TaskStatusClaimed {
				return errors.New("奖励已领取")
			}
			return errors.New("任务尚未完成")
		}
		if def.Period == model.TaskPeriodDaily && def.DailyLimit > 0 && p.ClaimCount >= def.DailyLimit {
			return errors.New("今日领取次数已达上限")
		}
		reward := def.RewardPoints
		now := common.LocalTime(time.Now())
		nextStatus := model.TaskStatusClaimed
		nextProgress := p.Progress
		nextClaim := p.ClaimCount + 1
		if def.Period == model.TaskPeriodDaily {
			if def.DailyLimit == 0 || nextClaim < def.DailyLimit {
				nextStatus = model.TaskStatusOngoing
				nextProgress = 0
			}
		}
		if _, err := session.ExecCtx(ctx,
			"UPDATE user_task_progress SET claim_count=?, claimed_at=?, status=?, progress=? WHERE id=?",
			nextClaim, &now, nextStatus, nextProgress, p.ID,
		); err != nil {
			return err
		}
		if reward <= 0 {
			up, err := r.ensurePointsTx(ctx, session, userID)
			if err != nil {
				return err
			}
			points = up
			return nil
		}
		up, err := r.addPointsTx(ctx, session, userID, reward, model.PointChangeTaskClaim, "任务奖励："+def.Title, "task", def.ID)
		if err != nil {
			return err
		}
		points = up
		return nil
	})
	return points, err
}

func (r *TaskRepository) ensurePointsTx(ctx context.Context, session sqlx.Session, userID uint64) (*model.UserPoints, error) {
	var up model.UserPoints
	err := session.QueryRowCtx(ctx, &up,
		"SELECT "+userPointsColumns+" FROM user_points WHERE user_id=? FOR UPDATE", userID,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		_, err = session.ExecCtx(ctx,
			"INSERT INTO user_points (user_id, points) VALUES (?, 0)", userID,
		)
		if err != nil {
			return nil, err
		}
		err = session.QueryRowCtx(ctx, &up,
			"SELECT "+userPointsColumns+" FROM user_points WHERE user_id=? FOR UPDATE", userID,
		)
	}
	if err != nil {
		return nil, err
	}
	return &up, nil
}

func (r *TaskRepository) addPointsTx(ctx context.Context, session sqlx.Session, userID uint64, delta int, changeType, remark, refType string, refID uint64) (*model.UserPoints, error) {
	up, err := r.ensurePointsTx(ctx, session, userID)
	if err != nil {
		return nil, err
	}
	up.Points += int64(delta)
	if up.Points < 0 {
		return nil, errors.New("积分不足")
	}
	if _, err := session.ExecCtx(ctx,
		"UPDATE user_points SET points=? WHERE user_id=?", up.Points, userID,
	); err != nil {
		return nil, err
	}
	if _, err := session.ExecCtx(ctx,
		`INSERT INTO user_point_logs (user_id, change_type, delta, points_after, remark, ref_type, ref_id)
		 VALUES (?,?,?,?,?,?,?)`,
		userID, changeType, delta, up.Points, remark, refType, refID,
	); err != nil {
		return nil, err
	}
	return up, nil
}

func (r *TaskRepository) DeductPoints(ctx context.Context, userID uint64, points int, changeType, remark, refType string, refID uint64) (*model.UserPoints, error) {
	if userID == 0 || points <= 0 {
		return nil, errors.New("参数无效")
	}
	if changeType == "" {
		changeType = model.PointChangeMallExchange
	}
	var out *model.UserPoints
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		n, err := countQuery(ctx, session,
			"SELECT COUNT(*) FROM user_point_logs WHERE ref_type=? AND ref_id=? AND change_type=?",
			refType, refID, changeType,
		)
		if err != nil {
			return err
		}
		if n > 0 {
			up, err := r.ensurePointsTx(ctx, session, userID)
			if err != nil {
				return err
			}
			out = up
			return nil
		}
		up, err := r.addPointsTx(ctx, session, userID, -points, changeType, remark, refType, refID)
		if err != nil {
			return err
		}
		out = up
		return nil
	})
	return out, err
}

func (r *TaskRepository) RefundPoints(ctx context.Context, userID uint64, points int, changeType, remark, refType string, refID uint64) (*model.UserPoints, error) {
	if userID == 0 || points <= 0 {
		return nil, errors.New("参数无效")
	}
	if changeType == "" {
		changeType = model.PointChangeMallRefund
	}
	var out *model.UserPoints
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		n, err := countQuery(ctx, session,
			"SELECT COUNT(*) FROM user_point_logs WHERE ref_type=? AND ref_id=? AND change_type=?",
			refType, refID, changeType,
		)
		if err != nil {
			return err
		}
		if n > 0 {
			up, err := r.ensurePointsTx(ctx, session, userID)
			if err != nil {
				return err
			}
			out = up
			return nil
		}
		up, err := r.addPointsTx(ctx, session, userID, points, changeType, remark, refType, refID)
		if err != nil {
			return err
		}
		out = up
		return nil
	})
	return out, err
}

func (r *TaskRepository) GetPoints(ctx context.Context, userID uint64) (*model.UserPoints, error) {
	var up model.UserPoints
	err := r.conn.QueryRowCtx(ctx, &up,
		"SELECT "+userPointsColumns+" FROM user_points WHERE user_id=? LIMIT 1", userID,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		up = model.UserPoints{UserID: userID, Points: 0}
		_, _ = r.conn.ExecCtx(ctx,
			"INSERT IGNORE INTO user_points (user_id, points) VALUES (?, 0)", userID,
		)
		return &up, nil
	}
	if err != nil {
		return nil, err
	}
	return &up, nil
}

func (r *TaskRepository) ListPointLogs(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserPointLog, int64, error) {
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_point_logs WHERE user_id=?", userID,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.UserPointLog
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+userPointLogColumns+" FROM user_point_logs WHERE user_id=? ORDER BY id DESC LIMIT ? OFFSET ?",
		userID, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *TaskRepository) GetUserBrief(ctx context.Context, userID uint64) (nickname, avatar string, err error) {
	var brief struct {
		Nickname string `db:"nickname"`
		Avatar   string `db:"avatar"`
	}
	err = r.conn.QueryRowCtx(ctx, &brief,
		"SELECT nickname, avatar FROM users WHERE id=? AND deleted_at IS NULL LIMIT 1", userID,
	)
	if err != nil {
		return "", "", err
	}
	return brief.Nickname, brief.Avatar, nil
}

func FormatDedupeKey(prefix string, id uint64) string {
	if id == 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", prefix, id)
}
