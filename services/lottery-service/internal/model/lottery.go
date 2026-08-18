package model

import (
	"database/sql"
	"time"
)

const (
	ActivityStatusDraft   = 0
	ActivityStatusOnline  = 1
	ActivityStatusOffline = 2

	PrizeTypePoints   = "points"
	PrizeTypeThanks   = "thanks"
	PrizeTypePhysical = "physical"

	RecordStatusPending = "pending"
	RecordStatusDone    = "done"
	RecordStatusFailed  = "failed"

	FulfillNone        = "none"
	FulfillNeedAddress = "need_address"
	FulfillPending     = "pending"
	FulfillShipped     = "shipped"
)

type LotteryActivity struct {
	ID         uint64     `db:"id"`
	Title      string     `db:"title"`
	Status     int        `db:"status"`
	CostPoints int        `db:"cost_points"`
	DailyLimit int        `db:"daily_limit"`
	StartAt    *time.Time `db:"start_at"`
	EndAt      *time.Time `db:"end_at"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

type LotteryPrize struct {
	ID           uint64    `db:"id"`
	ActivityID   uint64    `db:"activity_id"`
	Slot         int       `db:"slot"`
	Name         string    `db:"name"`
	CoverURL     string    `db:"cover_url"`
	PrizeType    string    `db:"prize_type"`
	PointsAmount int       `db:"points_amount"`
	Weight       int       `db:"weight"`
	Stock        int       `db:"stock"`
	StockStrict  int       `db:"stock_strict"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type LotteryDrawRecord struct {
	ID              uint64       `db:"id"`
	UserID          uint64       `db:"user_id"`
	ActivityID      uint64       `db:"activity_id"`
	PrizeID         uint64       `db:"prize_id"`
	Slot            int          `db:"slot"`
	PrizeType       string       `db:"prize_type"`
	PrizeName       string       `db:"prize_name"`
	PointsAmount    int          `db:"points_amount"`
	CostPoints      int          `db:"cost_points"`
	Status          string       `db:"status"`
	FulfillStatus   string       `db:"fulfill_status"`
	AddressID       uint64       `db:"address_id"`
	ReceiverName    string       `db:"receiver_name"`
	ReceiverPhone   string       `db:"receiver_phone"`
	ReceiverAddress string       `db:"receiver_address"`
	ShipCompany     string       `db:"ship_company"`
	ShipNo          string       `db:"ship_no"`
	ShippedAt       sql.NullTime `db:"shipped_at"`
	CreatedAt       time.Time    `db:"created_at"`
}
