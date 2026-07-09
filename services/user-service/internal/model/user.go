package model

import "mymall/common"

type BaseModel struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt common.LocalTime `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

type User struct {
	BaseModel
	Mobile        string           `gorm:"column:mobile;type:char(11);not null" json:"mobile"`
	Password      string           `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Nickname      string           `gorm:"column:nickname;type:varchar(50);not null" json:"nickname"`
	Avatar        string           `gorm:"column:avatar;type:varchar(255);default:''" json:"avatar"`
	Gender        int              `gorm:"column:gender;type:tinyint;default:0" json:"gender"`
	Status        int              `gorm:"column:status;type:tinyint;default:1" json:"status"`
	LastLoginTime common.LocalTime `gorm:"column:last_login_time;type:timestamp;default:null" json:"last_login_time"`
	DeletedAt     common.LocalTime `gorm:"column:deleted_at;type:timestamp;default:null;index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}
