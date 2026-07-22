package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductTagsModel = (*customProductTagsModel)(nil)

type (
	// ProductTagsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductTagsModel.
	ProductTagsModel interface {
		productTagsModel
		withSession(session sqlx.Session) ProductTagsModel
	}

	customProductTagsModel struct {
		*defaultProductTagsModel
	}
)

// NewProductTagsModel returns a model for the database table.
func NewProductTagsModel(conn sqlx.SqlConn) ProductTagsModel {
	return &customProductTagsModel{
		defaultProductTagsModel: newProductTagsModel(conn),
	}
}

func (m *customProductTagsModel) withSession(session sqlx.Session) ProductTagsModel {
	return NewProductTagsModel(sqlx.NewSqlConnFromSession(session))
}
