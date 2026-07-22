package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PointsExchangeOrdersModel = (*customPointsExchangeOrdersModel)(nil)

type (
	// PointsExchangeOrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPointsExchangeOrdersModel.
	PointsExchangeOrdersModel interface {
		pointsExchangeOrdersModel
		withSession(session sqlx.Session) PointsExchangeOrdersModel
	}

	customPointsExchangeOrdersModel struct {
		*defaultPointsExchangeOrdersModel
	}
)

// NewPointsExchangeOrdersModel returns a model for the database table.
func NewPointsExchangeOrdersModel(conn sqlx.SqlConn) PointsExchangeOrdersModel {
	return &customPointsExchangeOrdersModel{
		defaultPointsExchangeOrdersModel: newPointsExchangeOrdersModel(conn),
	}
}

func (m *customPointsExchangeOrdersModel) withSession(session sqlx.Session) PointsExchangeOrdersModel {
	return NewPointsExchangeOrdersModel(sqlx.NewSqlConnFromSession(session))
}
