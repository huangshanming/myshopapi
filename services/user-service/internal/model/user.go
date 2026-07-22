package model

import "mymall/common"

type BaseModel struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	CreatedAt common.LocalTime `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" db:"updated_at" json:"updated_at"`
}

type User struct {
	BaseModel
	Mobile        string           `gorm:"column:mobile;type:char(11);not null" db:"mobile" json:"mobile"`
	Password      string           `gorm:"column:password;type:varchar(255);not null" db:"password" json:"-"`
	Nickname      string           `gorm:"column:nickname;type:varchar(50);not null" db:"nickname" json:"nickname"`
	Avatar        string           `gorm:"column:avatar;type:varchar(255);default:''" db:"avatar" json:"avatar"`
	Gender        int              `gorm:"column:gender;type:tinyint;default:0" db:"gender" json:"gender"`
	Status        int              `gorm:"column:status;type:tinyint;default:1" db:"status" json:"status"`
	Role          string           `gorm:"column:role;type:varchar(32);default:user" db:"role" json:"role"`
	LastLoginTime common.LocalTime `gorm:"column:last_login_time;type:timestamp;default:null" db:"last_login_time" json:"last_login_time"`
	DeletedAt     common.LocalTime `gorm:"column:deleted_at;type:timestamp;default:null;index" db:"deleted_at" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}
