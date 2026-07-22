package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ OrderAfterSalesModel = (*customOrderAfterSalesModel)(nil)

type (
	// OrderAfterSalesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderAfterSalesModel.
	OrderAfterSalesModel interface {
		orderAfterSalesModel
		withSession(session sqlx.Session) OrderAfterSalesModel
	}

	customOrderAfterSalesModel struct {
		*defaultOrderAfterSalesModel
	}
)

// NewOrderAfterSalesModel returns a model for the database table.
func NewOrderAfterSalesModel(conn sqlx.SqlConn) OrderAfterSalesModel {
	return &customOrderAfterSalesModel{
		defaultOrderAfterSalesModel: newOrderAfterSalesModel(conn),
	}
}

func (m *customOrderAfterSalesModel) withSession(session sqlx.Session) OrderAfterSalesModel {
	return NewOrderAfterSalesModel(sqlx.NewSqlConnFromSession(session))
}
