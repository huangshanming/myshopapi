package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShopMenusModel = (*customShopMenusModel)(nil)

type (
	// ShopMenusModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShopMenusModel.
	ShopMenusModel interface {
		shopMenusModel
		withSession(session sqlx.Session) ShopMenusModel
	}

	customShopMenusModel struct {
		*defaultShopMenusModel
	}
)

// NewShopMenusModel returns a model for the database table.
func NewShopMenusModel(conn sqlx.SqlConn) ShopMenusModel {
	return &customShopMenusModel{
		defaultShopMenusModel: newShopMenusModel(conn),
	}
}

func (m *customShopMenusModel) withSession(session sqlx.Session) ShopMenusModel {
	return NewShopMenusModel(sqlx.NewSqlConnFromSession(session))
}
