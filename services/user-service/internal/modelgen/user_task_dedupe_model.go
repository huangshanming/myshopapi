package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserTaskDedupeModel = (*customUserTaskDedupeModel)(nil)

type (
	// UserTaskDedupeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTaskDedupeModel.
	UserTaskDedupeModel interface {
		userTaskDedupeModel
		withSession(session sqlx.Session) UserTaskDedupeModel
	}

	customUserTaskDedupeModel struct {
		*defaultUserTaskDedupeModel
	}
)

// NewUserTaskDedupeModel returns a model for the database table.
func NewUserTaskDedupeModel(conn sqlx.SqlConn) UserTaskDedupeModel {
	return &customUserTaskDedupeModel{
		defaultUserTaskDedupeModel: newUserTaskDedupeModel(conn),
	}
}

func (m *customUserTaskDedupeModel) withSession(session sqlx.Session) UserTaskDedupeModel {
	return NewUserTaskDedupeModel(sqlx.NewSqlConnFromSession(session))
}
