package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserWalletLogsModel = (*customUserWalletLogsModel)(nil)

type (
	// UserWalletLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserWalletLogsModel.
	UserWalletLogsModel interface {
		userWalletLogsModel
		withSession(session sqlx.Session) UserWalletLogsModel
	}

	customUserWalletLogsModel struct {
		*defaultUserWalletLogsModel
	}
)

// NewUserWalletLogsModel returns a model for the database table.
func NewUserWalletLogsModel(conn sqlx.SqlConn) UserWalletLogsModel {
	return &customUserWalletLogsModel{
		defaultUserWalletLogsModel: newUserWalletLogsModel(conn),
	}
}

func (m *customUserWalletLogsModel) withSession(session sqlx.Session) UserWalletLogsModel {
	return NewUserWalletLogsModel(sqlx.NewSqlConnFromSession(session))
}
