package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductOpLogsModel = (*customProductOpLogsModel)(nil)

type (
	// ProductOpLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductOpLogsModel.
	ProductOpLogsModel interface {
		productOpLogsModel
		withSession(session sqlx.Session) ProductOpLogsModel
	}

	customProductOpLogsModel struct {
		*defaultProductOpLogsModel
	}
)

// NewProductOpLogsModel returns a model for the database table.
func NewProductOpLogsModel(conn sqlx.SqlConn) ProductOpLogsModel {
	return &customProductOpLogsModel{
		defaultProductOpLogsModel: newProductOpLogsModel(conn),
	}
}

func (m *customProductOpLogsModel) withSession(session sqlx.Session) ProductOpLogsModel {
	return NewProductOpLogsModel(sqlx.NewSqlConnFromSession(session))
}
