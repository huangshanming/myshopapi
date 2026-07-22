package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LogisticsCompaniesModel = (*customLogisticsCompaniesModel)(nil)

type (
	// LogisticsCompaniesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLogisticsCompaniesModel.
	LogisticsCompaniesModel interface {
		logisticsCompaniesModel
		withSession(session sqlx.Session) LogisticsCompaniesModel
	}

	customLogisticsCompaniesModel struct {
		*defaultLogisticsCompaniesModel
	}
)

// NewLogisticsCompaniesModel returns a model for the database table.
func NewLogisticsCompaniesModel(conn sqlx.SqlConn) LogisticsCompaniesModel {
	return &customLogisticsCompaniesModel{
		defaultLogisticsCompaniesModel: newLogisticsCompaniesModel(conn),
	}
}

func (m *customLogisticsCompaniesModel) withSession(session sqlx.Session) LogisticsCompaniesModel {
	return NewLogisticsCompaniesModel(sqlx.NewSqlConnFromSession(session))
}
