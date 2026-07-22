package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductCategoriesModel = (*customProductCategoriesModel)(nil)

type (
	// ProductCategoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductCategoriesModel.
	ProductCategoriesModel interface {
		productCategoriesModel
		withSession(session sqlx.Session) ProductCategoriesModel
	}

	customProductCategoriesModel struct {
		*defaultProductCategoriesModel
	}
)

// NewProductCategoriesModel returns a model for the database table.
func NewProductCategoriesModel(conn sqlx.SqlConn) ProductCategoriesModel {
	return &customProductCategoriesModel{
		defaultProductCategoriesModel: newProductCategoriesModel(conn),
	}
}

func (m *customProductCategoriesModel) withSession(session sqlx.Session) ProductCategoriesModel {
	return NewProductCategoriesModel(sqlx.NewSqlConnFromSession(session))
}
