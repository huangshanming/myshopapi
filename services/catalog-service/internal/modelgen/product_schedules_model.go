package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductSchedulesModel = (*customProductSchedulesModel)(nil)

type (
	// ProductSchedulesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductSchedulesModel.
	ProductSchedulesModel interface {
		productSchedulesModel
		withSession(session sqlx.Session) ProductSchedulesModel
	}

	customProductSchedulesModel struct {
		*defaultProductSchedulesModel
	}
)

// NewProductSchedulesModel returns a model for the database table.
func NewProductSchedulesModel(conn sqlx.SqlConn) ProductSchedulesModel {
	return &customProductSchedulesModel{
		defaultProductSchedulesModel: newProductSchedulesModel(conn),
	}
}

func (m *customProductSchedulesModel) withSession(session sqlx.Session) ProductSchedulesModel {
	return NewProductSchedulesModel(sqlx.NewSqlConnFromSession(session))
}
