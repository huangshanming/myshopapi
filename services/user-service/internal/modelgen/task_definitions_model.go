package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TaskDefinitionsModel = (*customTaskDefinitionsModel)(nil)

type (
	// TaskDefinitionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskDefinitionsModel.
	TaskDefinitionsModel interface {
		taskDefinitionsModel
		withSession(session sqlx.Session) TaskDefinitionsModel
	}

	customTaskDefinitionsModel struct {
		*defaultTaskDefinitionsModel
	}
)

// NewTaskDefinitionsModel returns a model for the database table.
func NewTaskDefinitionsModel(conn sqlx.SqlConn) TaskDefinitionsModel {
	return &customTaskDefinitionsModel{
		defaultTaskDefinitionsModel: newTaskDefinitionsModel(conn),
	}
}

func (m *customTaskDefinitionsModel) withSession(session sqlx.Session) TaskDefinitionsModel {
	return NewTaskDefinitionsModel(sqlx.NewSqlConnFromSession(session))
}
