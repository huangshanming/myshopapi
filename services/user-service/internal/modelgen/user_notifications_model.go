package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserNotificationsModel = (*customUserNotificationsModel)(nil)

type (
	// UserNotificationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserNotificationsModel.
	UserNotificationsModel interface {
		userNotificationsModel
		withSession(session sqlx.Session) UserNotificationsModel
	}

	customUserNotificationsModel struct {
		*defaultUserNotificationsModel
	}
)

// NewUserNotificationsModel returns a model for the database table.
func NewUserNotificationsModel(conn sqlx.SqlConn) UserNotificationsModel {
	return &customUserNotificationsModel{
		defaultUserNotificationsModel: newUserNotificationsModel(conn),
	}
}

func (m *customUserNotificationsModel) withSession(session sqlx.Session) UserNotificationsModel {
	return NewUserNotificationsModel(sqlx.NewSqlConnFromSession(session))
}
