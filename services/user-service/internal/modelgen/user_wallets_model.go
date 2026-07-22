package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserWalletsModel = (*customUserWalletsModel)(nil)

type (
	// UserWalletsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserWalletsModel.
	UserWalletsModel interface {
		userWalletsModel
		withSession(session sqlx.Session) UserWalletsModel
	}

	customUserWalletsModel struct {
		*defaultUserWalletsModel
	}
)

// NewUserWalletsModel returns a model for the database table.
func NewUserWalletsModel(conn sqlx.SqlConn) UserWalletsModel {
	return &customUserWalletsModel{
		defaultUserWalletsModel: newUserWalletsModel(conn),
	}
}

func (m *customUserWalletsModel) withSession(session sqlx.Session) UserWalletsModel {
	return NewUserWalletsModel(sqlx.NewSqlConnFromSession(session))
}
