package model

import "mymall/common"

const (
	RegionLevelProvince = 1
	RegionLevelCity     = 2
	RegionLevelDistrict = 3
)

type Region struct {
	ID         uint64           `gorm:"column:id;primaryKey;autoIncrement" db:"id" json:"id"`
	Code       string           `gorm:"column:code;type:varchar(12);uniqueIndex" db:"code" json:"code"`
	Name       string           `gorm:"column:name;type:varchar(64)" db:"name" json:"name"`
	ParentCode string           `gorm:"column:parent_code;type:varchar(12);index" db:"parent_code" json:"parent_code"`
	Level      int              `gorm:"column:level" db:"level" json:"level"`
	Sort       int              `gorm:"column:sort" db:"sort" json:"sort"`
	CreatedAt  common.LocalTime `gorm:"column:created_at" db:"created_at" json:"created_at"`
}

func (Region) TableName() string { return "regions" }

type RegionTreeNode struct {
	Code     string           `json:"code"`
	Name     string           `json:"name"`
	Level    int              `json:"level"`
	Children []RegionTreeNode `json:"children,omitempty"`
}
