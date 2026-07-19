package model

import "mymall/common"

const (
	MsgTypeAnnounce = "announce"
	MsgTypeOrder    = "order"
	MsgTypeSystem   = "system"

	LinkTypeNone  = "none"
	LinkTypeOrder = "order"
	LinkTypeURL   = "url"

	SenderAdmin  = "admin"
	SenderSystem = "system"

	NotifyTargetAll   = "all"
	NotifyTargetUsers = "users"
)

type UserNotification struct {
	ID         uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     uint64           `gorm:"column:user_id;not null;index" json:"user_id"`
	Title      string           `gorm:"column:title;type:varchar(200)" json:"title"`
	Content    string           `gorm:"column:content;type:varchar(1000)" json:"content"`
	MsgType    string           `gorm:"column:msg_type;type:varchar(16)" json:"msg_type"`
	LinkType   string           `gorm:"column:link_type;type:varchar(16)" json:"link_type"`
	LinkID     uint64           `gorm:"column:link_id" json:"link_id"`
	Extra      string           `gorm:"column:extra;type:json" json:"extra,omitempty"`
	IsRead     int8             `gorm:"column:is_read" json:"is_read"`
	SenderType string           `gorm:"column:sender_type;type:varchar(16)" json:"sender_type"`
	SenderID   uint64           `gorm:"column:sender_id" json:"sender_id"`
	BatchID    uint64           `gorm:"column:batch_id" json:"batch_id"`
	CreatedAt  common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (UserNotification) TableName() string { return "user_notifications" }

type UserNotificationBatch struct {
	ID            uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title         string           `gorm:"column:title;type:varchar(200)" json:"title"`
	Content       string           `gorm:"column:content;type:varchar(1000)" json:"content"`
	Target        string           `gorm:"column:target;type:varchar(16)" json:"target"`
	UserCount     int              `gorm:"column:user_count" json:"user_count"`
	SuccessCount  int              `gorm:"column:success_count" json:"success_count"`
	LinkType      string           `gorm:"column:link_type;type:varchar(16)" json:"link_type"`
	LinkID        uint64           `gorm:"column:link_id" json:"link_id"`
	SenderID      uint64           `gorm:"column:sender_id" json:"sender_id"`
	CreatedAt     common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (UserNotificationBatch) TableName() string { return "user_notification_batches" }
