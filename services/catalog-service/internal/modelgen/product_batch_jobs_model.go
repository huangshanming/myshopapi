package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductBatchJobsModel = (*customProductBatchJobsModel)(nil)

type (
	// ProductBatchJobsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductBatchJobsModel.
	ProductBatchJobsModel interface {
		productBatchJobsModel
		withSession(session sqlx.Session) ProductBatchJobsModel
	}

	customProductBatchJobsModel struct {
		*defaultProductBatchJobsModel
	}
)

// NewProductBatchJobsModel returns a model for the database table.
func NewProductBatchJobsModel(conn sqlx.SqlConn) ProductBatchJobsModel {
	return &customProductBatchJobsModel{
		defaultProductBatchJobsModel: newProductBatchJobsModel(conn),
	}
}

func (m *customProductBatchJobsModel) withSession(session sqlx.Session) ProductBatchJobsModel {
	return NewProductBatchJobsModel(sqlx.NewSqlConnFromSession(session))
}
