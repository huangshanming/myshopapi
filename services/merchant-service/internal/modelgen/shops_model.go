package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShopsModel = (*customShopsModel)(nil)

type (
	// ShopsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShopsModel.
	ShopsModel interface {
		shopsModel
		withSession(session sqlx.Session) ShopsModel
	}

	customShopsModel struct {
		*defaultShopsModel
	}
)

// NewShopsModel returns a model for the database table.
func NewShopsModel(conn sqlx.SqlConn) ShopsModel {
	return &customShopsModel{
		defaultShopsModel: newShopsModel(conn),
	}
}

func (m *customShopsModel) withSession(session sqlx.Session) ShopsModel {
	return NewShopsModel(sqlx.NewSqlConnFromSession(session))
}
