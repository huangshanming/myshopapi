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
	ID                 uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	Name               string           `gorm:"column:name;type:varchar(100);not null" db:"name" json:"name"`
	Logo               string           `gorm:"column:logo;type:varchar(500)" db:"logo" json:"logo"`
	ContactName        string           `gorm:"column:contact_name;type:varchar(50)" db:"contact_name" json:"contact_name"`
	ContactPhone       string           `gorm:"column:contact_phone;type:char(11)" db:"contact_phone" json:"contact_phone"`
	Description        string           `gorm:"column:description;type:varchar(500)" db:"description" json:"description"`
	Category           string           `gorm:"column:category;type:varchar(50)" db:"category" json:"category"`
	Province           string           `gorm:"column:province;type:varchar(50)" db:"province" json:"province"`
	City               string           `gorm:"column:city;type:varchar(50)" db:"city" json:"city"`
	District           string           `gorm:"column:district;type:varchar(50)" db:"district" json:"district"`
	Address            string           `gorm:"column:address;type:varchar(255)" db:"address" json:"address"`
	BusinessLicenseNo  string           `gorm:"column:business_license_no;type:varchar(64)" db:"business_license_no" json:"business_license_no"`
	LegalPerson        string           `gorm:"column:legal_person;type:varchar(50)" db:"legal_person" json:"legal_person"`
	LicenseImage       string           `gorm:"column:license_image;type:varchar(500)" db:"license_image" json:"license_image"`
	StorefrontImage    string           `gorm:"column:storefront_image;type:varchar(500)" db:"storefront_image" json:"storefront_image"`
	OwnerUserID        uint64           `gorm:"column:owner_user_id;not null;index" db:"owner_user_id" json:"owner_user_id"`
	Status             string           `gorm:"column:status;type:enum('pending','approved','rejected','disabled');default:pending" db:"status" json:"status"`
	RejectReason       string           `gorm:"column:reject_reason;type:varchar(255)" db:"reject_reason" json:"reject_reason"`
	CreatedAt          common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt          common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (Shop) TableName() string { return "shops" }

type ShopApplication struct {
	ID                uint64            `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	UserID            uint64            `gorm:"column:user_id;not null;index" db:"user_id" json:"user_id"`
	ShopName          string            `gorm:"column:shop_name;type:varchar(100);not null" db:"shop_name" json:"shop_name"`
	ContactName       string            `gorm:"column:contact_name;type:varchar(50);not null" db:"contact_name" json:"contact_name"`
	ContactPhone      string            `gorm:"column:contact_phone;type:char(11);not null" db:"contact_phone" json:"contact_phone"`
	Description       string            `gorm:"column:description;type:varchar(500)" db:"description" json:"description"`
	Category          string            `gorm:"column:category;type:varchar(50)" db:"category" json:"category"`
	Province          string            `gorm:"column:province;type:varchar(50)" db:"province" json:"province"`
	City              string            `gorm:"column:city;type:varchar(50)" db:"city" json:"city"`
	District          string            `gorm:"column:district;type:varchar(50)" db:"district" json:"district"`
	Address           string            `gorm:"column:address;type:varchar(255)" db:"address" json:"address"`
	BusinessLicenseNo string            `gorm:"column:business_license_no;type:varchar(64)" db:"business_license_no" json:"business_license_no"`
	LegalPerson       string            `gorm:"column:legal_person;type:varchar(50)" db:"legal_person" json:"legal_person"`
	LicenseImage      string            `gorm:"column:license_image;type:varchar(500)" db:"license_image" json:"license_image"`
	StorefrontImage   string            `gorm:"column:storefront_image;type:varchar(500)" db:"storefront_image" json:"storefront_image"`
	Status            string            `gorm:"column:status;type:enum('pending','approved','rejected');default:pending" db:"status" json:"status"`
	RejectReason      string            `gorm:"column:reject_reason;type:varchar(255)" db:"reject_reason" json:"reject_reason"`
	// Non-pointer: go-zero sqlx pre-allocates *uint64 so NULL cannot scan; use IFNULL in SELECT.
	ReviewedBy        uint64            `gorm:"column:reviewed_by" db:"reviewed_by" json:"reviewed_by,omitempty"`
	ReviewedAt        *common.LocalTime `gorm:"column:reviewed_at" db:"reviewed_at" json:"reviewed_at,omitempty"`
	ShopID            uint64            `gorm:"column:shop_id" db:"shop_id" json:"shop_id,omitempty"`
	CreatedAt         common.LocalTime  `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt         common.LocalTime  `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (ShopApplication) TableName() string { return "shop_applications" }

type ShopMember struct {
	ID         uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	ShopID     uint64           `gorm:"column:shop_id;not null;index" db:"shop_id" json:"shop_id"`
	UserID     uint64           `gorm:"column:user_id;not null;index" db:"user_id" json:"user_id"`
	MemberRole string           `gorm:"column:member_role;type:enum('owner','staff');default:owner" db:"member_role" json:"member_role"`
	CreatedAt  common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt  common.LocalTime `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (ShopMember) TableName() string { return "shop_members" }
