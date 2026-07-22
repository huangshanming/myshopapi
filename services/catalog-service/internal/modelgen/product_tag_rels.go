package modelgen

type ProductTagRel struct {
	ProductID uint64 `db:"product_id"`
	TagID     uint64 `db:"tag_id"`
}

func (ProductTagRel) TableName() string { return "product_tag_rels" }
