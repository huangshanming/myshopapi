package model

const (
	RegionLevelProvince = 1
	RegionLevelCity     = 2
	RegionLevelDistrict = 3
)

type RegionTreeNode struct {
	Code     string           `json:"code"`
	Name     string           `json:"name"`
	Level    int              `json:"level"`
	Children []RegionTreeNode `json:"children,omitempty"`
}
