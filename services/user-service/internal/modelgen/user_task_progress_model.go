package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserTaskProgressModel = (*customUserTaskProgressModel)(nil)

type (
	// UserTaskProgressModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTaskProgressModel.
	UserTaskProgressModel interface {
		userTaskProgressModel
		withSession(session sqlx.Session) UserTaskProgressModel
	}

	customUserTaskProgressModel struct {
		*defaultUserTaskProgressModel
	}
)

// NewUserTaskProgressModel returns a model for the database table.
func NewUserTaskProgressModel(conn sqlx.SqlConn) UserTaskProgressModel {
	return &customUserTaskProgressModel{
		defaultUserTaskProgressModel: newUserTaskProgressModel(conn),
	}
}

func (m *customUserTaskProgressModel) withSession(session sqlx.Session) UserTaskProgressModel {
	return NewUserTaskProgressModel(sqlx.NewSqlConnFromSession(session))
}
