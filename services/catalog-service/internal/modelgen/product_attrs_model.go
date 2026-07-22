package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductAttrsModel = (*customProductAttrsModel)(nil)

type (
	// ProductAttrsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductAttrsModel.
	ProductAttrsModel interface {
		productAttrsModel
		withSession(session sqlx.Session) ProductAttrsModel
	}

	customProductAttrsModel struct {
		*defaultProductAttrsModel
	}
)

// NewProductAttrsModel returns a model for the database table.
func NewProductAttrsModel(conn sqlx.SqlConn) ProductAttrsModel {
	return &customProductAttrsModel{
		defaultProductAttrsModel: newProductAttrsModel(conn),
	}
}

func (m *customProductAttrsModel) withSession(session sqlx.Session) ProductAttrsModel {
	return NewProductAttrsModel(sqlx.NewSqlConnFromSession(session))
}
