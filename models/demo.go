package model

// User 对应数据库 user 表（GORM 映射配置）
type Demo struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`                // 主键、自增
	Name string `json:"name" gorm:"column:name;type:varchar(50);not null"` // 字段名、类型、非空
}

// TableName 自定义表名（GORM 默认会将结构体名转为小写复数，如 User -> users，此方法可覆盖默认规则）
func (u *Demo) TableName() string {
	return "demo" // 对应数据库中的 user 表
}
