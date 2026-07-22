package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserNotificationBatchesModel = (*customUserNotificationBatchesModel)(nil)

type (
	// UserNotificationBatchesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserNotificationBatchesModel.
	UserNotificationBatchesModel interface {
		userNotificationBatchesModel
		withSession(session sqlx.Session) UserNotificationBatchesModel
	}

	customUserNotificationBatchesModel struct {
		*defaultUserNotificationBatchesModel
	}
)

// NewUserNotificationBatchesModel returns a model for the database table.
func NewUserNotificationBatchesModel(conn sqlx.SqlConn) UserNotificationBatchesModel {
	return &customUserNotificationBatchesModel{
		defaultUserNotificationBatchesModel: newUserNotificationBatchesModel(conn),
	}
}

func (m *customUserNotificationBatchesModel) withSession(session sqlx.Session) UserNotificationBatchesModel {
	return NewUserNotificationBatchesModel(sqlx.NewSqlConnFromSession(session))
}
