package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductSkusModel = (*customProductSkusModel)(nil)

type (
	// ProductSkusModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductSkusModel.
	ProductSkusModel interface {
		productSkusModel
		withSession(session sqlx.Session) ProductSkusModel
	}

	customProductSkusModel struct {
		*defaultProductSkusModel
	}
)

// NewProductSkusModel returns a model for the database table.
func NewProductSkusModel(conn sqlx.SqlConn) ProductSkusModel {
	return &customProductSkusModel{
		defaultProductSkusModel: newProductSkusModel(conn),
	}
}

func (m *customProductSkusModel) withSession(session sqlx.Session) ProductSkusModel {
	return NewProductSkusModel(sqlx.NewSqlConnFromSession(session))
}
