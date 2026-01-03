package model

type ProductCategory struct {
	BaseModel
	ParentId     uint64 `gorm:"column:parent_id;type:int;default:0;comment:父分类ID" json:"parent_id"`
	Name         string `gorm:"column:name;type:varchar(100);not null;comment:分类名称" json:"name"`
	Icon         string `gorm:"column:icon;type:varchar(200);default:null;comment:分类图标" json:"icon,omitempty"`
	Image        string `gorm:"column:image;type:varchar(500);default:null;comment:分类图片" json:"image,omitempty"`
	Description  string `gorm:"column:description;type:varchar(500);default:null;comment:分类描述" json:"description,omitempty"`
	SortOrder    int    `gorm:"column:sort_order;type:int;default:0;comment:排序" json:"sort_order"`
	Level        int    `gorm:"column:level;type:int;default:1;comment:分类层级" json:"level"`
	IsShow       bool   `gorm:"column:is_show;type:tinyint(1);default:1;comment:是否显示" json:"is_show"`
	ProductCount int    `gorm:"column:product_count;type:int;default:0;comment:商品数量" json:"product_count"`
}

func (p *ProductCategory) TableName() string {
	return "product_categories"
}
