package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShopWalletLogsModel = (*customShopWalletLogsModel)(nil)

type (
	// ShopWalletLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShopWalletLogsModel.
	ShopWalletLogsModel interface {
		shopWalletLogsModel
		withSession(session sqlx.Session) ShopWalletLogsModel
	}

	customShopWalletLogsModel struct {
		*defaultShopWalletLogsModel
	}
)

// NewShopWalletLogsModel returns a model for the database table.
func NewShopWalletLogsModel(conn sqlx.SqlConn) ShopWalletLogsModel {
	return &customShopWalletLogsModel{
		defaultShopWalletLogsModel: newShopWalletLogsModel(conn),
	}
}

func (m *customShopWalletLogsModel) withSession(session sqlx.Session) ShopWalletLogsModel {
	return NewShopWalletLogsModel(sqlx.NewSqlConnFromSession(session))
}
