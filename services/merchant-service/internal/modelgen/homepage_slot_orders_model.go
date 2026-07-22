package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomepageSlotOrdersModel = (*customHomepageSlotOrdersModel)(nil)

type (
	// HomepageSlotOrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomepageSlotOrdersModel.
	HomepageSlotOrdersModel interface {
		homepageSlotOrdersModel
		withSession(session sqlx.Session) HomepageSlotOrdersModel
	}

	customHomepageSlotOrdersModel struct {
		*defaultHomepageSlotOrdersModel
	}
)

// NewHomepageSlotOrdersModel returns a model for the database table.
func NewHomepageSlotOrdersModel(conn sqlx.SqlConn) HomepageSlotOrdersModel {
	return &customHomepageSlotOrdersModel{
		defaultHomepageSlotOrdersModel: newHomepageSlotOrdersModel(conn),
	}
}

func (m *customHomepageSlotOrdersModel) withSession(session sqlx.Session) HomepageSlotOrdersModel {
	return NewHomepageSlotOrdersModel(sqlx.NewSqlConnFromSession(session))
}
