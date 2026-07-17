package logic

import (
	"mymall/services/catalog-service/internal/notify/repository"
	"mymall/services/catalog-service/internal/svc"
)

type NotificationLogic struct {
	svcCtx *svc.ServiceContext
}

func NewNotificationLogic(svcCtx *svc.ServiceContext) *NotificationLogic {
	return &NotificationLogic{svcCtx: svcCtx}
}

func (l *NotificationLogic) List(f repository.NotificationListFilter) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.Notifications.List(f)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (l *NotificationLogic) UnreadCount(shopID uint64) (map[string]interface{}, error) {
	cnt, err := l.svcCtx.Notifications.UnreadCount(shopID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"count": cnt}, nil
}

func (l *NotificationLogic) MarkRead(id, shopID uint64) error {
	return l.svcCtx.Notifications.MarkRead(id, shopID)
}

func (l *NotificationLogic) MarkAllRead(shopID uint64) error {
	return l.svcCtx.Notifications.MarkAllRead(shopID)
}
