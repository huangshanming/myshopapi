package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductFavoritesModel = (*customProductFavoritesModel)(nil)

type (
	// ProductFavoritesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductFavoritesModel.
	ProductFavoritesModel interface {
		productFavoritesModel
		withSession(session sqlx.Session) ProductFavoritesModel
	}

	customProductFavoritesModel struct {
		*defaultProductFavoritesModel
	}
)

// NewProductFavoritesModel returns a model for the database table.
func NewProductFavoritesModel(conn sqlx.SqlConn) ProductFavoritesModel {
	return &customProductFavoritesModel{
		defaultProductFavoritesModel: newProductFavoritesModel(conn),
	}
}

func (m *customProductFavoritesModel) withSession(session sqlx.Session) ProductFavoritesModel {
	return NewProductFavoritesModel(sqlx.NewSqlConnFromSession(session))
}
