package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomepageSlotPackagesModel = (*customHomepageSlotPackagesModel)(nil)

type (
	// HomepageSlotPackagesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomepageSlotPackagesModel.
	HomepageSlotPackagesModel interface {
		homepageSlotPackagesModel
		withSession(session sqlx.Session) HomepageSlotPackagesModel
	}

	customHomepageSlotPackagesModel struct {
		*defaultHomepageSlotPackagesModel
	}
)

// NewHomepageSlotPackagesModel returns a model for the database table.
func NewHomepageSlotPackagesModel(conn sqlx.SqlConn) HomepageSlotPackagesModel {
	return &customHomepageSlotPackagesModel{
		defaultHomepageSlotPackagesModel: newHomepageSlotPackagesModel(conn),
	}
}

func (m *customHomepageSlotPackagesModel) withSession(session sqlx.Session) HomepageSlotPackagesModel {
	return NewHomepageSlotPackagesModel(sqlx.NewSqlConnFromSession(session))
}
