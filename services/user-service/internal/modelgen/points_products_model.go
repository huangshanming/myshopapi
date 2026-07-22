package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PointsProductsModel = (*customPointsProductsModel)(nil)

type (
	// PointsProductsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPointsProductsModel.
	PointsProductsModel interface {
		pointsProductsModel
		withSession(session sqlx.Session) PointsProductsModel
	}

	customPointsProductsModel struct {
		*defaultPointsProductsModel
	}
)

// NewPointsProductsModel returns a model for the database table.
func NewPointsProductsModel(conn sqlx.SqlConn) PointsProductsModel {
	return &customPointsProductsModel{
		defaultPointsProductsModel: newPointsProductsModel(conn),
	}
}

func (m *customPointsProductsModel) withSession(session sqlx.Session) PointsProductsModel {
	return NewPointsProductsModel(sqlx.NewSqlConnFromSession(session))
}
