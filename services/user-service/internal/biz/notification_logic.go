package biz

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/repository"
)

type NotifyCreateReq struct {
	UserID   uint64 `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	MsgType  string `json:"msg_type"`
	LinkType string `json:"link_type"`
	LinkID   uint64 `json:"link_id"`
	Extra    string `json:"extra"`
}

type AdminSendReq struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Target   string   `json:"target"` // all | users
	UserIDs  []uint64 `json:"user_ids"`
	LinkType string   `json:"link_type"`
	LinkID   uint64   `json:"link_id"`
}

func (l *UserLogic) CreateNotification(ctx context.Context, req NotifyCreateReq) (*model.UserNotification, error) {
	if req.UserID == 0 {
		return nil, errors.New("缺少用户")
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("标题不能为空")
	}
	msgType := req.MsgType
	if msgType == "" {
		msgType = model.MsgTypeSystem
	}
	linkType := req.LinkType
	if linkType == "" {
		linkType = model.LinkTypeNone
	}
	extra := strings.TrimSpace(req.Extra)
	if extra == "" {
		extra = "{}"
	}
	n := &model.UserNotification{
		UserID:     req.UserID,
		Title:      title,
		Content:    strings.TrimSpace(req.Content),
		MsgType:    msgType,
		LinkType:   linkType,
		LinkID:     req.LinkID,
		Extra:      extra,
		IsRead:     0,
		SenderType: model.SenderSystem,
		SenderID:   0,
	}
	if err := l.svcCtx.Repo.CreateNotification(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (l *UserLogic) ListMyNotifications(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserNotification, int64, error) {
	return l.svcCtx.Repo.ListNotifications(ctx, userID, page, pageSize)
}

func (l *UserLogic) UnreadCount(ctx context.Context, userID uint64) (int64, error) {
	return l.svcCtx.Repo.UnreadNotificationCount(ctx, userID)
}

func (l *UserLogic) MarkRead(ctx context.Context, userID, id uint64) error {
	return l.svcCtx.Repo.MarkNotificationRead(ctx, userID, id)
}

func (l *UserLogic) MarkAllRead(ctx context.Context, userID uint64) error {
	return l.svcCtx.Repo.MarkAllNotificationsRead(ctx, userID)
}

func (l *UserLogic) AdminSend(ctx context.Context, adminID uint64, req AdminSendReq) (*model.UserNotificationBatch, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("标题不能为空")
	}
	target := req.Target
	if target == "" {
		target = model.NotifyTargetUsers
	}
	linkType := req.LinkType
	if linkType == "" {
		linkType = model.LinkTypeNone
	}
	var userIDs []uint64
	if target == model.NotifyTargetAll {
		const batchSize = 500
		offset := 0
		for {
			ids, err := l.svcCtx.Repo.ListActiveUserIDs(ctx, offset, batchSize)
			if err != nil {
				return nil, err
			}
			if len(ids) == 0 {
				break
			}
			userIDs = append(userIDs, ids...)
			if len(userIDs) >= 5000 {
				userIDs = userIDs[:5000]
				break
			}
			if len(ids) < batchSize {
				break
			}
			offset += batchSize
		}
	} else {
		seen := map[uint64]struct{}{}
		for _, id := range req.UserIDs {
			if id == 0 {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			userIDs = append(userIDs, id)
			if len(userIDs) >= 5000 {
				break
			}
		}
	}
	if len(userIDs) == 0 {
		return nil, errors.New("没有可发送的用户")
	}
	batch := &model.UserNotificationBatch{
		Title: title, Content: strings.TrimSpace(req.Content), Target: target,
		UserCount: len(userIDs), LinkType: linkType, LinkID: req.LinkID, SenderID: adminID,
	}
	if err := l.svcCtx.Repo.CreateNotificationBatch(ctx, batch); err != nil {
		return nil, err
	}
	list := make([]model.UserNotification, 0, len(userIDs))
	content := strings.TrimSpace(req.Content)
	for _, uid := range userIDs {
		list = append(list, model.UserNotification{
			UserID: uid, Title: title, Content: content, MsgType: model.MsgTypeAnnounce,
			LinkType: linkType, LinkID: req.LinkID, Extra: "{}", IsRead: 0,
			SenderType: model.SenderAdmin, SenderID: adminID, BatchID: batch.ID,
		})
	}
	if err := l.svcCtx.Repo.CreateNotifications(ctx, list); err != nil {
		return nil, err
	}
	batch.SuccessCount = len(list)
	_ = l.svcCtx.Repo.UpdateNotificationBatchSuccess(ctx, batch.ID, batch.SuccessCount)
	return batch, nil
}

func (l *UserLogic) ListSendBatches(ctx context.Context, page, pageSize int) ([]model.UserNotificationBatch, int64, error) {
	return l.svcCtx.Repo.ListNotificationBatches(ctx, page, pageSize)
}

func (l *UserLogic) GetSendBatch(ctx context.Context, id uint64) (*model.UserNotificationBatch, error) {
	return l.svcCtx.Repo.GetNotificationBatch(ctx, id)
}

func (l *UserLogic) ListBatchRecipients(ctx context.Context, batchID uint64, page, pageSize int) ([]repository.BatchRecipientRow, int64, error) {
	return l.svcCtx.Repo.ListBatchRecipients(ctx, batchID, page, pageSize)
}

// ExtraJSON helper for order service callers via internal API
func ExtraJSON(v map[string]interface{}) string {
	if v == nil {
		return ""
	}
	b, _ := json.Marshal(v)
	return string(b)
}
