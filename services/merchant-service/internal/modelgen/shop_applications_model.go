package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShopApplicationsModel = (*customShopApplicationsModel)(nil)

type (
	// ShopApplicationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShopApplicationsModel.
	ShopApplicationsModel interface {
		shopApplicationsModel
		withSession(session sqlx.Session) ShopApplicationsModel
	}

	customShopApplicationsModel struct {
		*defaultShopApplicationsModel
	}
)

// NewShopApplicationsModel returns a model for the database table.
func NewShopApplicationsModel(conn sqlx.SqlConn) ShopApplicationsModel {
	return &customShopApplicationsModel{
		defaultShopApplicationsModel: newShopApplicationsModel(conn),
	}
}

func (m *customShopApplicationsModel) withSession(session sqlx.Session) ShopApplicationsModel {
	return NewShopApplicationsModel(sqlx.NewSqlConnFromSession(session))
}
