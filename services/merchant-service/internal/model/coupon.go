package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"mymall/common"
)

const (
	CouponIssuerPlatform = "platform"
	CouponIssuerShop     = "shop"

	CouponTypeFullReduce  = "full_reduce"
	CouponTypeNoThreshold = "no_threshold"
	CouponTypeCategory    = "category"
	CouponTypeProduct     = "product"
	CouponTypeDiscount    = "discount"

	CouponScopeAll      = "all"
	CouponScopeCategory = "category"
	CouponScopeProduct  = "product"

	CouponValidFixed    = "fixed"
	CouponValidRelative = "relative"

	CouponStatusDraft   = "draft"
	CouponStatusOn      = "on"
	CouponStatusOff     = "off"
	CouponStatusExpired = "expired"

	CouponIdentityAll = "all"
	CouponIdentityNew = "new"
	CouponIdentityOld = "old"

	UserCouponUnused  = "unused"
	UserCouponLocked  = "locked"
	UserCouponUsed    = "used"
	UserCouponExpired = "expired"

	CouponSourceDirect    = "direct"
	CouponSourceOrderGift = "order_gift"
	CouponSourcePopup     = "popup"
	CouponSourceTargeted  = "targeted"

	CouponActionRedeem = "redeem"
	CouponActionUnlock = "unlock"
	CouponActionReturn = "return"

	CouponChannelDirect    = "direct"
	CouponChannelOrderGift = "order_gift"
	CouponChannelPopup     = "popup"
	CouponChannelTargeted  = "targeted"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return errors.New("invalid channels type")
	}
	if len(b) == 0 {
		*s = StringSlice{}
		return nil
	}
	return json.Unmarshal(b, s)
}

type Coupon struct {
	ID                uint64            `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	Name              string            `gorm:"column:name;type:varchar(100)" db:"name" json:"name"`
	IssuerType        string            `gorm:"column:issuer_type;type:varchar(16)" db:"issuer_type" json:"issuer_type"`
	ShopID            uint64            `gorm:"column:shop_id" db:"shop_id" json:"shop_id"`
	CouponType        string            `gorm:"column:coupon_type;type:varchar(32)" db:"coupon_type" json:"coupon_type"`
	ThresholdAmount   float64           `gorm:"column:threshold_amount;type:decimal(12,2)" db:"threshold_amount" json:"threshold_amount"`
	DiscountAmount    float64           `gorm:"column:discount_amount;type:decimal(12,2)" db:"discount_amount" json:"discount_amount"`
	DiscountRate      float64           `gorm:"column:discount_rate;type:decimal(6,4)" db:"discount_rate" json:"discount_rate"`
	MaxDiscountAmount float64           `gorm:"column:max_discount_amount;type:decimal(12,2)" db:"max_discount_amount" json:"max_discount_amount"`
	ScopeType         string            `gorm:"column:scope_type;type:varchar(16)" db:"scope_type" json:"scope_type"`
	TotalCount        int               `gorm:"column:total_count" db:"total_count" json:"total_count"`
	ClaimedCount      int               `gorm:"column:claimed_count" db:"claimed_count" json:"claimed_count"`
	PerUserLimit      int               `gorm:"column:per_user_limit" db:"per_user_limit" json:"per_user_limit"`
	ValidType         string            `gorm:"column:valid_type;type:varchar(16)" db:"valid_type" json:"valid_type"`
	ValidStart        *common.LocalTime `gorm:"column:valid_start" db:"valid_start" json:"valid_start,omitempty"`
	ValidEnd          *common.LocalTime `gorm:"column:valid_end" db:"valid_end" json:"valid_end,omitempty"`
	ValidDays         int               `gorm:"column:valid_days" db:"valid_days" json:"valid_days"`
	Stackable         int8              `gorm:"column:stackable" db:"stackable" json:"stackable"`
	UserIdentity      string            `gorm:"column:user_identity;type:varchar(16)" db:"user_identity" json:"user_identity"`
	Channels          StringSlice       `gorm:"column:channels;type:json" db:"channels" json:"channels"`
	Status            string            `gorm:"column:status;type:varchar(16)" db:"status" json:"status"`
	Remark            string            `gorm:"column:remark;type:varchar(255)" db:"remark" json:"remark"`
	CreatedBy         uint64            `gorm:"column:created_by" db:"created_by" json:"created_by"`
	CreatedAt         common.LocalTime  `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt         common.LocalTime  `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
	Scopes            []CouponScope     `gorm:"-" db:"-" json:"scopes,omitempty"`
	DisplayStatus     string            `gorm:"-" db:"-" json:"display_status,omitempty"`
	ClaimedByMe       bool              `gorm:"-" db:"-" json:"claimed_by_me,omitempty"`
	Remaining         *int              `gorm:"-" db:"-" json:"remaining,omitempty"`
}

func (Coupon) TableName() string { return "coupons" }

type CouponScope struct {
	ID       uint64 `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	CouponID uint64 `gorm:"column:coupon_id;index" db:"coupon_id" json:"coupon_id"`
	RefType  string `gorm:"column:ref_type;type:varchar(16)" db:"ref_type" json:"ref_type"`
	RefID    uint64 `gorm:"column:ref_id" db:"ref_id" json:"ref_id"`
}

func (CouponScope) TableName() string { return "coupon_scopes" }

type UserCoupon struct {
	ID             uint64            `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	CouponID       uint64            `gorm:"column:coupon_id;index" db:"coupon_id" json:"coupon_id"`
	UserID         uint64            `gorm:"column:user_id;index" db:"user_id" json:"user_id"`
	ShopID         uint64            `gorm:"column:shop_id" db:"shop_id" json:"shop_id"`
	Status         string            `gorm:"column:status;type:varchar(16)" db:"status" json:"status"`
	Source         string            `gorm:"column:source;type:varchar(16)" db:"source" json:"source"`
	ValidStart     common.LocalTime  `gorm:"column:valid_start" db:"valid_start" json:"valid_start"`
	ValidEnd       common.LocalTime  `gorm:"column:valid_end" db:"valid_end" json:"valid_end"`
	OrderID        uint64            `gorm:"column:order_id" db:"order_id" json:"order_id"`
	LockedAt       *common.LocalTime `gorm:"column:locked_at" db:"locked_at" json:"locked_at,omitempty"`
	UsedAt         *common.LocalTime `gorm:"column:used_at" db:"used_at" json:"used_at,omitempty"`
	ClaimBatchNo   string            `gorm:"column:claim_batch_no;type:varchar(64)" db:"claim_batch_no" json:"claim_batch_no"`
	DiscountAmount float64           `gorm:"column:discount_amount;type:decimal(12,2)" db:"discount_amount" json:"discount_amount"`
	CreatedAt      common.LocalTime  `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt      common.LocalTime  `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
	Coupon         *Coupon           `gorm:"-" db:"-" json:"coupon,omitempty"`
	CouponName     string            `gorm:"-" db:"-" json:"coupon_name,omitempty"`
}

func (UserCoupon) TableName() string { return "user_coupons" }

type CouponGrant struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	CouponID     uint64           `gorm:"column:coupon_id" db:"coupon_id" json:"coupon_id"`
	OperatorID   uint64           `gorm:"column:operator_id" db:"operator_id" json:"operator_id"`
	IssuerType   string           `gorm:"column:issuer_type" db:"issuer_type" json:"issuer_type"`
	ShopID       uint64           `gorm:"column:shop_id" db:"shop_id" json:"shop_id"`
	UserCount    int              `gorm:"column:user_count" db:"user_count" json:"user_count"`
	SuccessCount int              `gorm:"column:success_count" db:"success_count" json:"success_count"`
	BatchNo      string           `gorm:"column:batch_no" db:"batch_no" json:"batch_no"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
}

func (CouponGrant) TableName() string { return "coupon_grants" }

type CouponRedeemLog struct {
	ID             uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	UserCouponID   uint64           `gorm:"column:user_coupon_id" db:"user_coupon_id" json:"user_coupon_id"`
	CouponID       uint64           `gorm:"column:coupon_id" db:"coupon_id" json:"coupon_id"`
	UserID         uint64           `gorm:"column:user_id" db:"user_id" json:"user_id"`
	OrderID        uint64           `gorm:"column:order_id" db:"order_id" json:"order_id"`
	ShopID         uint64           `gorm:"column:shop_id" db:"shop_id" json:"shop_id"`
	DiscountAmount float64          `gorm:"column:discount_amount;type:decimal(12,2)" db:"discount_amount" json:"discount_amount"`
	Action         string           `gorm:"column:action;type:varchar(16)" db:"action" json:"action"`
	CreatedAt      common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
}

func (CouponRedeemLog) TableName() string { return "coupon_redeem_logs" }
