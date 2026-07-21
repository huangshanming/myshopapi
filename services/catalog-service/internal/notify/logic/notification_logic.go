package logic

import (
	"context"
	"mymall/services/catalog-service/internal/notify/repository"
	"mymall/services/catalog-service/internal/svc"
)

type NotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationLogic {
	return &NotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotificationLogic) List(f repository.NotificationListFilter) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.Notifications.List(l.ctx, f)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (l *NotificationLogic) UnreadCount(shopID uint64) (map[string]interface{}, error) {
	cnt, err := l.svcCtx.Notifications.UnreadCount(l.ctx, shopID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"count": cnt}, nil
}

func (l *NotificationLogic) MarkRead(id, shopID uint64) error {
	return l.svcCtx.Notifications.MarkRead(l.ctx, id, shopID)
}

func (l *NotificationLogic) MarkAllRead(shopID uint64) error {
	return l.svcCtx.Notifications.MarkAllRead(l.ctx, shopID)
}
