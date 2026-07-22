package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShopRolesModel = (*customShopRolesModel)(nil)

type (
	// ShopRolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShopRolesModel.
	ShopRolesModel interface {
		shopRolesModel
		withSession(session sqlx.Session) ShopRolesModel
	}

	customShopRolesModel struct {
		*defaultShopRolesModel
	}
)

// NewShopRolesModel returns a model for the database table.
func NewShopRolesModel(conn sqlx.SqlConn) ShopRolesModel {
	return &customShopRolesModel{
		defaultShopRolesModel: newShopRolesModel(conn),
	}
}

func (m *customShopRolesModel) withSession(session sqlx.Session) ShopRolesModel {
	return NewShopRolesModel(sqlx.NewSqlConnFromSession(session))
}
