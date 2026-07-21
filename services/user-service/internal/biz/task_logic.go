package biz

import (
	"context"
	"errors"
	"strings"

	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/repository"
	"mymall/services/user-service/internal/svc"
)

type TaskLogic struct {
	svcCtx *svc.ServiceContext
}

func NewTaskLogic(svcCtx *svc.ServiceContext) *TaskLogic {
	return &TaskLogic{svcCtx: svcCtx}
}

type TaskItemVO struct {
	Code         string `json:"code"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	Period       string `json:"period"`
	RewardPoints int    `json:"reward_points"`
	TargetCount  int    `json:"target_count"`
	DailyLimit   int    `json:"daily_limit"`
	Progress     int    `json:"progress"`
	ClaimCount   int    `json:"claim_count"`
	Status       string `json:"status"`
	Enabled      int8   `json:"enabled"`
}

func (l *TaskLogic) ListUserTasks(ctx context.Context, userID uint64) ([]TaskItemVO, error) {
	defs, err := l.svcCtx.Tasks.ListDefinitions(ctx, false)
	if err != nil {
		return nil, err
	}
	out := make([]TaskItemVO, 0, len(defs))
	for _, d := range defs {
		p, err := l.svcCtx.Tasks.GetOrCreateProgress(ctx, userID, &d)
		if err != nil {
			return nil, err
		}
		out = append(out, TaskItemVO{
			Code: d.Code, Title: d.Title, Description: d.Description, Icon: d.Icon,
			Period: d.Period, RewardPoints: d.RewardPoints, TargetCount: d.TargetCount,
			DailyLimit: d.DailyLimit, Progress: p.Progress, ClaimCount: p.ClaimCount,
			Status: p.Status, Enabled: d.Enabled,
		})
	}
	return out, nil
}

func (l *TaskLogic) AdminListTasks(ctx context.Context) ([]model.TaskDefinition, error) {
	return l.svcCtx.Tasks.ListDefinitions(ctx, true)
}

type UpdateTaskReq struct {
	Title        *string `json:"title"`
	Description  *string `json:"description"`
	Icon         *string `json:"icon"`
	Enabled      *int8   `json:"enabled"`
	RewardPoints *int    `json:"reward_points"`
	TargetCount  *int    `json:"target_count"`
	DailyLimit   *int    `json:"daily_limit"`
	Sort         *int    `json:"sort"`
	RulesJSON    *string `json:"rules_json"`
}

func (l *TaskLogic) AdminUpdateTask(ctx context.Context, id uint64, req UpdateTaskReq) (*model.TaskDefinition, error) {
	if _, err := l.svcCtx.Tasks.GetDefinitionByID(ctx, id); err != nil {
		return nil, errors.New("任务不存在")
	}
	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		updates["description"] = strings.TrimSpace(*req.Description)
	}
	if req.Icon != nil {
		updates["icon"] = strings.TrimSpace(*req.Icon)
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.RewardPoints != nil {
		if *req.RewardPoints < 0 {
			return nil, errors.New("积分不能为负")
		}
		updates["reward_points"] = *req.RewardPoints
	}
	if req.TargetCount != nil {
		if *req.TargetCount < 1 {
			return nil, errors.New("目标次数至少为 1")
		}
		updates["target_count"] = *req.TargetCount
	}
	if req.DailyLimit != nil {
		if *req.DailyLimit < 0 {
			return nil, errors.New("每日上限无效")
		}
		updates["daily_limit"] = *req.DailyLimit
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.RulesJSON != nil {
		updates["rules_json"] = *req.RulesJSON
	}
	if len(updates) == 0 {
		return l.svcCtx.Tasks.GetDefinitionByID(ctx, id)
	}
	if err := l.svcCtx.Tasks.UpdateDefinition(ctx, id, updates); err != nil {
		return nil, err
	}
	return l.svcCtx.Tasks.GetDefinitionByID(ctx, id)
}

type TaskEventReq struct {
	UserID   uint64 `json:"user_id"`
	TaskCode string `json:"task_code"`
	Delta    int    `json:"delta"`
	RefType  string `json:"ref_type"`
	RefID    uint64 `json:"ref_id"`
}

func (l *TaskLogic) HandleEvent(ctx context.Context, req TaskEventReq) error {
	if req.UserID == 0 || strings.TrimSpace(req.TaskCode) == "" {
		return errors.New("参数无效")
	}
	def, err := l.svcCtx.Tasks.GetDefinition(ctx, req.TaskCode)
	if err != nil {
		return nil // 未知任务忽略
	}
	if def.Enabled != 1 {
		return nil
	}
	// 特殊：完善资料需校验
	if req.TaskCode == "first_profile" {
		nick, avatar, err := l.svcCtx.Tasks.GetUserBrief(ctx, req.UserID)
		if err != nil || strings.TrimSpace(nick) == "" || strings.TrimSpace(avatar) == "" {
			return nil
		}
	}
	refKey := repository.FormatDedupeKey(req.RefType, req.RefID)
	if refKey == "" && req.RefType != "" {
		refKey = req.RefType
	}
	_, err = l.svcCtx.Tasks.ApplyEvent(ctx, req.UserID, def, req.Delta, refKey)
	return err
}

func (l *TaskLogic) Checkin(ctx context.Context, userID uint64) (*model.UserTaskProgress, error) {
	def, err := l.svcCtx.Tasks.GetDefinition(ctx, "daily_checkin")
	if err != nil {
		return nil, errors.New("签到任务未配置")
	}
	if def.Enabled != 1 {
		return nil, errors.New("签到任务已关闭")
	}
	p, err := l.svcCtx.Tasks.GetOrCreateProgress(ctx, userID, def)
	if err != nil {
		return nil, err
	}
	if p.Status == model.TaskStatusClaimable {
		return p, nil
	}
	if p.Status == model.TaskStatusClaimed && def.DailyLimit > 0 && p.ClaimCount >= def.DailyLimit {
		return nil, errors.New("今日已签到")
	}
	return l.svcCtx.Tasks.ApplyEvent(ctx, userID, def, 1, "checkin:"+repository.TodayBizDate())
}

func (l *TaskLogic) Claim(ctx context.Context, userID uint64, code string) (*model.UserPoints, error) {
	def, err := l.svcCtx.Tasks.GetDefinition(ctx, code)
	if err != nil {
		return nil, errors.New("任务不存在")
	}
	if def.Enabled != 1 {
		return nil, errors.New("任务已关闭")
	}
	return l.svcCtx.Tasks.Claim(ctx, userID, def)
}

func (l *TaskLogic) GetPoints(ctx context.Context, userID uint64) (*model.UserPoints, error) {
	return l.svcCtx.Tasks.GetPoints(ctx, userID)
}

type PointsLedgerReq struct {
	UserID     uint64 `json:"user_id"`
	Points     int    `json:"points"`
	ChangeType string `json:"change_type"`
	Remark     string `json:"remark"`
	RefType    string `json:"ref_type"`
	RefID      uint64 `json:"ref_id"`
}

func (l *TaskLogic) DeductPoints(ctx context.Context, req PointsLedgerReq) (*model.UserPoints, error) {
	if req.UserID == 0 || req.Points <= 0 {
		return nil, errors.New("参数无效")
	}
	if strings.TrimSpace(req.RefType) == "" || req.RefID == 0 {
		return nil, errors.New("缺少业务单号")
	}
	return l.svcCtx.Tasks.DeductPoints(ctx, req.UserID, req.Points, req.ChangeType, req.Remark, req.RefType, req.RefID)
}

func (l *TaskLogic) RefundPoints(ctx context.Context, req PointsLedgerReq) (*model.UserPoints, error) {
	if req.UserID == 0 || req.Points <= 0 {
		return nil, errors.New("参数无效")
	}
	if strings.TrimSpace(req.RefType) == "" || req.RefID == 0 {
		return nil, errors.New("缺少业务单号")
	}
	return l.svcCtx.Tasks.RefundPoints(ctx, req.UserID, req.Points, req.ChangeType, req.Remark, req.RefType, req.RefID)
}

func (l *TaskLogic) ListPointLogs(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserPointLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return l.svcCtx.Tasks.ListPointLogs(ctx, userID, page, pageSize)
}
