package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShopMembersModel = (*customShopMembersModel)(nil)

type (
	// ShopMembersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShopMembersModel.
	ShopMembersModel interface {
		shopMembersModel
		withSession(session sqlx.Session) ShopMembersModel
	}

	customShopMembersModel struct {
		*defaultShopMembersModel
	}
)

// NewShopMembersModel returns a model for the database table.
func NewShopMembersModel(conn sqlx.SqlConn) ShopMembersModel {
	return &customShopMembersModel{
		defaultShopMembersModel: newShopMembersModel(conn),
	}
}

func (m *customShopMembersModel) withSession(session sqlx.Session) ShopMembersModel {
	return NewShopMembersModel(sqlx.NewSqlConnFromSession(session))
}
