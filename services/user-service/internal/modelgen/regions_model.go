package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RegionsModel = (*customRegionsModel)(nil)

type (
	// RegionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRegionsModel.
	RegionsModel interface {
		regionsModel
		withSession(session sqlx.Session) RegionsModel
	}

	customRegionsModel struct {
		*defaultRegionsModel
	}
)

// NewRegionsModel returns a model for the database table.
func NewRegionsModel(conn sqlx.SqlConn) RegionsModel {
	return &customRegionsModel{
		defaultRegionsModel: newRegionsModel(conn),
	}
}

func (m *customRegionsModel) withSession(session sqlx.Session) RegionsModel {
	return NewRegionsModel(sqlx.NewSqlConnFromSession(session))
}
