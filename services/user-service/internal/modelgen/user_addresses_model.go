package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserAddressesModel = (*customUserAddressesModel)(nil)

type (
	// UserAddressesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserAddressesModel.
	UserAddressesModel interface {
		userAddressesModel
		withSession(session sqlx.Session) UserAddressesModel
	}

	customUserAddressesModel struct {
		*defaultUserAddressesModel
	}
)

// NewUserAddressesModel returns a model for the database table.
func NewUserAddressesModel(conn sqlx.SqlConn) UserAddressesModel {
	return &customUserAddressesModel{
		defaultUserAddressesModel: newUserAddressesModel(conn),
	}
}

func (m *customUserAddressesModel) withSession(session sqlx.Session) UserAddressesModel {
	return NewUserAddressesModel(sqlx.NewSqlConnFromSession(session))
}
