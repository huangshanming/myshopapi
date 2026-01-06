package model

import (
	"mymall/common"
)

// Product 商品主表模型，对应数据库表products
// 字段类型严格匹配数据库表定义，标签包含GORM配置、JSON序列化配置及注释
type User struct {
	// 基础字段
	BaseModel
	Mobile        string           `gorm:"column:mobile;type:char(11);not null;comment:登录手机号" json:"mobile"`
	Password      string           `gorm:"column:password;type:varchar(255);not null;comment:登录密码" json:"password"`
	Nickname      string           `gorm:"column:nickname;type:varchar(50);not null;comment:用户昵称" json:"nickname"`
	Avatar        string           `gorm:"column:avatar;type:varchar(255);default:'';comment:用户头像URL" json:"avatar"`
	Gender        int              `gorm:"column:gender;type:tinyint;default:0;comment:性别：0-未知 1-男 2-女" json:"gender"`
	Status        int              `gorm:"column:status;type:tinyint;default:1;comment:账号状态：1-正常 0-禁用" json:"status"`
	LastLoginTime common.LocalTime `gorm:"column:last_login_time;type:timestamp;default:null;comment:最后登录时间" json:"last_login_time"`
	DeletedAt     common.LocalTime `gorm:"column:deleted_at;type:timestamp;default:null;index;comment:删除时间（软删除）" json:"deleted_at"`
}

// TableName 指定结构体对应的数据库表名
// GORM默认会将结构体名转为蛇形复数（如Product转为products），此处显式指定确保一致性
func (p *User) TableName() string {
	return "users"
}
