package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserPointLogsModel = (*customUserPointLogsModel)(nil)

type (
	// UserPointLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserPointLogsModel.
	UserPointLogsModel interface {
		userPointLogsModel
		withSession(session sqlx.Session) UserPointLogsModel
	}

	customUserPointLogsModel struct {
		*defaultUserPointLogsModel
	}
)

// NewUserPointLogsModel returns a model for the database table.
func NewUserPointLogsModel(conn sqlx.SqlConn) UserPointLogsModel {
	return &customUserPointLogsModel{
		defaultUserPointLogsModel: newUserPointLogsModel(conn),
	}
}

func (m *customUserPointLogsModel) withSession(session sqlx.Session) UserPointLogsModel {
	return NewUserPointLogsModel(sqlx.NewSqlConnFromSession(session))
}
