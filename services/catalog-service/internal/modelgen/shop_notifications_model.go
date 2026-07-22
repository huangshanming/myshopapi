package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShopNotificationsModel = (*customShopNotificationsModel)(nil)

type (
	// ShopNotificationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShopNotificationsModel.
	ShopNotificationsModel interface {
		shopNotificationsModel
		withSession(session sqlx.Session) ShopNotificationsModel
	}

	customShopNotificationsModel struct {
		*defaultShopNotificationsModel
	}
)

// NewShopNotificationsModel returns a model for the database table.
func NewShopNotificationsModel(conn sqlx.SqlConn) ShopNotificationsModel {
	return &customShopNotificationsModel{
		defaultShopNotificationsModel: newShopNotificationsModel(conn),
	}
}

func (m *customShopNotificationsModel) withSession(session sqlx.Session) ShopNotificationsModel {
	return NewShopNotificationsModel(sqlx.NewSqlConnFromSession(session))
}
