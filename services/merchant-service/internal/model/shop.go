package model

import "mymall/common"

const (
	ShopPending  = "pending"
	ShopApproved = "approved"
	ShopRejected = "rejected"
	ShopDisabled = "disabled"

	AppPending  = "pending"
	AppApproved = "approved"
	AppRejected = "rejected"

	MemberOwner = "owner"
	MemberStaff = "staff"
)

type Shop struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string           `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Logo         string           `gorm:"column:logo;type:varchar(500)" json:"logo"`
	ContactName  string           `gorm:"column:contact_name;type:varchar(50)" json:"contact_name"`
	ContactPhone string           `gorm:"column:contact_phone;type:char(11)" json:"contact_phone"`
	Description  string           `gorm:"column:description;type:varchar(500)" json:"description"`
	OwnerUserID  uint64           `gorm:"column:owner_user_id;not null;index" json:"owner_user_id"`
	Status       string           `gorm:"column:status;type:enum('pending','approved','rejected','disabled');default:pending" json:"status"`
	RejectReason string           `gorm:"column:reject_reason;type:varchar(255)" json:"reject_reason"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (Shop) TableName() string { return "shops" }

type ShopApplication struct {
	ID           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       uint64           `gorm:"column:user_id;not null;index" json:"user_id"`
	ShopName     string           `gorm:"column:shop_name;type:varchar(100);not null" json:"shop_name"`
	ContactName  string           `gorm:"column:contact_name;type:varchar(50);not null" json:"contact_name"`
	ContactPhone string           `gorm:"column:contact_phone;type:char(11);not null" json:"contact_phone"`
	Description  string           `gorm:"column:description;type:varchar(500)" json:"description"`
	Status       string           `gorm:"column:status;type:enum('pending','approved','rejected');default:pending" json:"status"`
	RejectReason string           `gorm:"column:reject_reason;type:varchar(255)" json:"reject_reason"`
	ReviewedBy   *uint64          `gorm:"column:reviewed_by" json:"reviewed_by,omitempty"`
	ReviewedAt   *common.LocalTime `gorm:"column:reviewed_at" json:"reviewed_at,omitempty"`
	ShopID       *uint64          `gorm:"column:shop_id" json:"shop_id,omitempty"`
	CreatedAt    common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ShopApplication) TableName() string { return "shop_applications" }

type ShopMember struct {
	ID         uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID     uint64           `gorm:"column:shop_id;not null;index" json:"shop_id"`
	UserID     uint64           `gorm:"column:user_id;not null;index" json:"user_id"`
	MemberRole string           `gorm:"column:member_role;type:enum('owner','staff');default:owner" json:"member_role"`
	CreatedAt  common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ShopMember) TableName() string { return "shop_members" }
