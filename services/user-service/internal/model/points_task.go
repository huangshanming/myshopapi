package model

import "mymall/common"

const (
	TaskPeriodDaily = "daily"
	TaskPeriodOnce  = "once"

	TaskStatusOngoing   = "ongoing"
	TaskStatusClaimable = "claimable"
	TaskStatusClaimed   = "claimed"

	PointChangeTaskClaim    = "task_claim"
	PointChangeAdminAdjust  = "admin_adjust"
	PointChangeMallExchange = "points_mall_exchange"
	PointChangeMallRefund   = "points_mall_refund"

	TaskBizDateOnce = "1970-01-01"
)

type UserPoints struct {
	UserID    uint64           `gorm:"column:user_id;primaryKey" db:"user_id" json:"user_id"`
	Points    int64            `gorm:"column:points" db:"points" json:"points"`
	CreatedAt common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (UserPoints) TableName() string { return "user_points" }

type UserPointLog struct {
	ID          uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	UserID      uint64           `gorm:"column:user_id;index" db:"user_id" json:"user_id"`
	ChangeType  string           `gorm:"column:change_type;type:varchar(32)" db:"change_type" json:"change_type"`
	Delta       int              `gorm:"column:delta" db:"delta" json:"delta"`
	PointsAfter int64            `gorm:"column:points_after" db:"points_after" json:"points_after"`
	Remark      string           `gorm:"column:remark;type:varchar(255)" db:"remark" json:"remark"`
	RefType     string           `gorm:"column:ref_type;type:varchar(32)" db:"ref_type" json:"ref_type"`
	RefID       uint64           `gorm:"column:ref_id" db:"ref_id" json:"ref_id"`
	CreatedAt   common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
}

func (UserPointLog) TableName() string { return "user_point_logs" }

type TaskDefinition struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	Code         string           `gorm:"column:code;type:varchar(64);uniqueIndex" db:"code" json:"code"`
	Title        string           `gorm:"column:title;type:varchar(100)" db:"title" json:"title"`
	Description  string           `gorm:"column:description;type:varchar(500)" db:"description" json:"description"`
	Icon         string           `gorm:"column:icon;type:varchar(64)" db:"icon" json:"icon"`
	Period       string           `gorm:"column:period;type:varchar(16)" db:"period" json:"period"`
	Enabled      int8             `gorm:"column:enabled" db:"enabled" json:"enabled"`
	RewardPoints int              `gorm:"column:reward_points" db:"reward_points" json:"reward_points"`
	TargetCount  int              `gorm:"column:target_count" db:"target_count" json:"target_count"`
	DailyLimit   int              `gorm:"column:daily_limit" db:"daily_limit" json:"daily_limit"`
	Sort         int              `gorm:"column:sort" db:"sort" json:"sort"`
	RulesJSON    string           `gorm:"column:rules_json;type:varchar(1000)" db:"rules_json" json:"rules_json"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (TaskDefinition) TableName() string { return "task_definitions" }

type UserTaskProgress struct {
	ID         uint64            `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	UserID     uint64            `gorm:"column:user_id" db:"user_id" json:"user_id"`
	TaskCode   string            `gorm:"column:task_code;type:varchar(64)" db:"task_code" json:"task_code"`
	BizDate    string            `gorm:"column:biz_date;type:date" db:"biz_date" json:"biz_date"`
	Progress   int               `gorm:"column:progress" db:"progress" json:"progress"`
	ClaimCount int               `gorm:"column:claim_count" db:"claim_count" json:"claim_count"`
	Status     string            `gorm:"column:status;type:varchar(16)" db:"status" json:"status"`
	ClaimedAt  *common.LocalTime `gorm:"column:claimed_at" db:"claimed_at" json:"claimed_at,omitempty"`
	CreatedAt  common.LocalTime  `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt  common.LocalTime  `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (UserTaskProgress) TableName() string { return "user_task_progress" }

type UserTaskDedupe struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	UserID    uint64           `gorm:"column:user_id" db:"user_id" json:"user_id"`
	TaskCode  string           `gorm:"column:task_code;type:varchar(64)" db:"task_code" json:"task_code"`
	BizDate   string           `gorm:"column:biz_date;type:date" db:"biz_date" json:"biz_date"`
	RefKey    string           `gorm:"column:ref_key;type:varchar(64)" db:"ref_key" json:"ref_key"`
	CreatedAt common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
}

func (UserTaskDedupe) TableName() string { return "user_task_dedupe" }
