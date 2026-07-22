package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductImagesModel = (*customProductImagesModel)(nil)

type (
	// ProductImagesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductImagesModel.
	ProductImagesModel interface {
		productImagesModel
		withSession(session sqlx.Session) ProductImagesModel
	}

	customProductImagesModel struct {
		*defaultProductImagesModel
	}
)

// NewProductImagesModel returns a model for the database table.
func NewProductImagesModel(conn sqlx.SqlConn) ProductImagesModel {
	return &customProductImagesModel{
		defaultProductImagesModel: newProductImagesModel(conn),
	}
}

func (m *customProductImagesModel) withSession(session sqlx.Session) ProductImagesModel {
	return NewProductImagesModel(sqlx.NewSqlConnFromSession(session))
}
