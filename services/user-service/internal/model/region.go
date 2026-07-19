package model

import "mymall/common"

const (
	RegionLevelProvince = 1
	RegionLevelCity     = 2
	RegionLevelDistrict = 3
)

type Region struct {
	ID         uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code       string           `gorm:"column:code;type:varchar(12);uniqueIndex" json:"code"`
	Name       string           `gorm:"column:name;type:varchar(64)" json:"name"`
	ParentCode string           `gorm:"column:parent_code;type:varchar(12);index" json:"parent_code"`
	Level      int              `gorm:"column:level" json:"level"`
	Sort       int              `gorm:"column:sort" json:"sort"`
	CreatedAt  common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (Region) TableName() string { return "regions" }

type RegionTreeNode struct {
	Code     string           `json:"code"`
	Name     string           `json:"name"`
	Level    int              `json:"level"`
	Children []RegionTreeNode `json:"children,omitempty"`
}
