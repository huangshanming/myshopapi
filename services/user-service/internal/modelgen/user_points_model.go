package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserPointsModel = (*customUserPointsModel)(nil)

type (
	// UserPointsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserPointsModel.
	UserPointsModel interface {
		userPointsModel
		withSession(session sqlx.Session) UserPointsModel
	}

	customUserPointsModel struct {
		*defaultUserPointsModel
	}
)

// NewUserPointsModel returns a model for the database table.
func NewUserPointsModel(conn sqlx.SqlConn) UserPointsModel {
	return &customUserPointsModel{
		defaultUserPointsModel: newUserPointsModel(conn),
	}
}

func (m *customUserPointsModel) withSession(session sqlx.Session) UserPointsModel {
	return NewUserPointsModel(sqlx.NewSqlConnFromSession(session))
}
