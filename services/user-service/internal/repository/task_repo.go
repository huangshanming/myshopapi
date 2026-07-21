package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mymall/common"
	"mymall/services/user-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) SeedIfEmpty(ctx context.Context) error {
	var n int64
	if err := r.db.WithContext(ctx).Model(&model.TaskDefinition{}).Count(&n).Error; err != nil {
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
	return r.db.WithContext(ctx).Create(&seeds).Error
}

func (r *TaskRepository) ListDefinitions(ctx context.Context, all bool) ([]model.TaskDefinition, error) {
	q := r.db.WithContext(ctx).Model(&model.TaskDefinition{}).Order("sort ASC, id ASC")
	if !all {
		q = q.Where("enabled = 1")
	}
	var list []model.TaskDefinition
	err := q.Find(&list).Error
	return list, err
}

func (r *TaskRepository) GetDefinition(ctx context.Context, code string) (*model.TaskDefinition, error) {
	var t model.TaskDefinition
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) GetDefinitionByID(ctx context.Context, id uint64) (*model.TaskDefinition, error) {
	var t model.TaskDefinition
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) UpdateDefinition(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.TaskDefinition{}).Where("id = ?", id).Updates(updates).Error
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
	err := r.db.WithContext(ctx).Where("user_id = ? AND task_code = ? AND biz_date = ?", userID, def.Code, biz).First(&p).Error
	if err == nil {
		return &p, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	p = model.UserTaskProgress{
		UserID: userID, TaskCode: def.Code, BizDate: biz,
		Progress: 0, ClaimCount: 0, Status: model.TaskStatusOngoing,
	}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&p).Error; err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Where("user_id = ? AND task_code = ? AND biz_date = ?", userID, def.Code, biz).First(&p).Error
	return &p, err
}

func (r *TaskRepository) TryDedupe(ctx context.Context, userID uint64, taskCode, bizDate, refKey string) (bool, error) {
	if refKey == "" {
		return true, nil
	}
	row := model.UserTaskDedupe{UserID: userID, TaskCode: taskCode, BizDate: bizDate, RefKey: refKey}
	res := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&row)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
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
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var p model.UserTaskProgress
		e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND task_code = ? AND biz_date = ?", userID, def.Code, biz).
			First(&p).Error
		if errors.Is(e, gorm.ErrRecordNotFound) {
			p = model.UserTaskProgress{
				UserID: userID, TaskCode: def.Code, BizDate: biz,
				Status: model.TaskStatusOngoing,
			}
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
			e = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", p.ID).First(&p).Error
		}
		if e != nil {
			return e
		}
		// 已领满（once 或 daily_limit）
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
			// 日任务领取后若还可再领一轮，重置进度
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
		if err := tx.Model(&p).Updates(map[string]interface{}{
			"progress": p.Progress, "status": p.Status,
		}).Error; err != nil {
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
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var p model.UserTaskProgress
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND task_code = ? AND biz_date = ?", userID, def.Code, biz).
			First(&p).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
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
		updates := map[string]interface{}{
			"claim_count": p.ClaimCount + 1,
			"claimed_at":  &now,
		}
		// 领取后：once 永久 claimed；daily 若还可再做则回到 ongoing
		nextStatus := model.TaskStatusClaimed
		if def.Period == model.TaskPeriodDaily {
			nextClaim := p.ClaimCount + 1
			if def.DailyLimit == 0 || nextClaim < def.DailyLimit {
				nextStatus = model.TaskStatusOngoing
				updates["progress"] = 0
			}
		}
		updates["status"] = nextStatus
		if err := tx.Model(&p).Updates(updates).Error; err != nil {
			return err
		}
		if reward <= 0 {
			up, err := r.ensurePointsTx(tx, userID)
			if err != nil {
				return err
			}
			points = up
			return nil
		}
		up, err := r.addPointsTx(tx, userID, reward, model.PointChangeTaskClaim, "任务奖励："+def.Title, "task", def.ID)
		if err != nil {
			return err
		}
		points = up
		return nil
	})
	return points, err
}

func (r *TaskRepository) ensurePointsTx(tx *gorm.DB, userID uint64) (*model.UserPoints, error) {
	var up model.UserPoints
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&up).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		up = model.UserPoints{UserID: userID, Points: 0}
		if err := tx.Create(&up).Error; err != nil {
			return nil, err
		}
		return &up, nil
	}
	if err != nil {
		return nil, err
	}
	return &up, nil
}

func (r *TaskRepository) addPointsTx(tx *gorm.DB, userID uint64, delta int, changeType, remark, refType string, refID uint64) (*model.UserPoints, error) {
	up, err := r.ensurePointsTx(tx, userID)
	if err != nil {
		return nil, err
	}
	up.Points += int64(delta)
	if up.Points < 0 {
		return nil, errors.New("积分不足")
	}
	if err := tx.Model(up).Update("points", up.Points).Error; err != nil {
		return nil, err
	}
	log := model.UserPointLog{
		UserID: userID, ChangeType: changeType, Delta: delta,
		PointsAfter: up.Points, Remark: remark, RefType: refType, RefID: refID,
	}
	if err := tx.Create(&log).Error; err != nil {
		return nil, err
	}
	return up, nil
}

// DeductPoints 扣减积分（幂等：同 ref_type+ref_id+change_type 已存在则直接成功）
func (r *TaskRepository) DeductPoints(ctx context.Context, userID uint64, points int, changeType, remark, refType string, refID uint64) (*model.UserPoints, error) {
	if userID == 0 || points <= 0 {
		return nil, errors.New("参数无效")
	}
	if changeType == "" {
		changeType = model.PointChangeMallExchange
	}
	var out *model.UserPoints
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var n int64
		if err := tx.Model(&model.UserPointLog{}).
			Where("ref_type = ? AND ref_id = ? AND change_type = ?", refType, refID, changeType).
			Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			up, err := r.ensurePointsTx(tx, userID)
			if err != nil {
				return err
			}
			out = up
			return nil
		}
		up, err := r.addPointsTx(tx, userID, -points, changeType, remark, refType, refID)
		if err != nil {
			return err
		}
		out = up
		return nil
	})
	return out, err
}

// RefundPoints 退回积分（幂等）
func (r *TaskRepository) RefundPoints(ctx context.Context, userID uint64, points int, changeType, remark, refType string, refID uint64) (*model.UserPoints, error) {
	if userID == 0 || points <= 0 {
		return nil, errors.New("参数无效")
	}
	if changeType == "" {
		changeType = model.PointChangeMallRefund
	}
	var out *model.UserPoints
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var n int64
		if err := tx.Model(&model.UserPointLog{}).
			Where("ref_type = ? AND ref_id = ? AND change_type = ?", refType, refID, changeType).
			Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			up, err := r.ensurePointsTx(tx, userID)
			if err != nil {
				return err
			}
			out = up
			return nil
		}
		up, err := r.addPointsTx(tx, userID, points, changeType, remark, refType, refID)
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
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&up).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		up = model.UserPoints{UserID: userID, Points: 0}
		_ = r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&up).Error
		return &up, nil
	}
	if err != nil {
		return nil, err
	}
	return &up, nil
}

func (r *TaskRepository) ListPointLogs(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserPointLog, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&model.UserPointLog{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.UserPointLog
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *TaskRepository) GetUserBrief(ctx context.Context, userID uint64) (nickname, avatar string, err error) {
	var u model.User
	err = r.db.WithContext(ctx).Select("id", "nickname", "avatar").Where("id = ?", userID).First(&u).Error
	if err != nil {
		return "", "", err
	}
	return u.Nickname, u.Avatar, nil
}

func FormatDedupeKey(prefix string, id uint64) string {
	if id == 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", prefix, id)
}
