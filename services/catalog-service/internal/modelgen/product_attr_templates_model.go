package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductAttrTemplatesModel = (*customProductAttrTemplatesModel)(nil)

type (
	// ProductAttrTemplatesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductAttrTemplatesModel.
	ProductAttrTemplatesModel interface {
		productAttrTemplatesModel
		withSession(session sqlx.Session) ProductAttrTemplatesModel
	}

	customProductAttrTemplatesModel struct {
		*defaultProductAttrTemplatesModel
	}
)

// NewProductAttrTemplatesModel returns a model for the database table.
func NewProductAttrTemplatesModel(conn sqlx.SqlConn) ProductAttrTemplatesModel {
	return &customProductAttrTemplatesModel{
		defaultProductAttrTemplatesModel: newProductAttrTemplatesModel(conn),
	}
}

func (m *customProductAttrTemplatesModel) withSession(session sqlx.Session) ProductAttrTemplatesModel {
	return NewProductAttrTemplatesModel(sqlx.NewSqlConnFromSession(session))
}
