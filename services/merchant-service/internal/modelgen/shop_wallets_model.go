package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShopWalletsModel = (*customShopWalletsModel)(nil)

type (
	// ShopWalletsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShopWalletsModel.
	ShopWalletsModel interface {
		shopWalletsModel
		withSession(session sqlx.Session) ShopWalletsModel
	}

	customShopWalletsModel struct {
		*defaultShopWalletsModel
	}
)

// NewShopWalletsModel returns a model for the database table.
func NewShopWalletsModel(conn sqlx.SqlConn) ShopWalletsModel {
	return &customShopWalletsModel{
		defaultShopWalletsModel: newShopWalletsModel(conn),
	}
}

func (m *customShopWalletsModel) withSession(session sqlx.Session) ShopWalletsModel {
	return NewShopWalletsModel(sqlx.NewSqlConnFromSession(session))
}
