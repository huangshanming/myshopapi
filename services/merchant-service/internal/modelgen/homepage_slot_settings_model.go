package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomepageSlotSettingsModel = (*customHomepageSlotSettingsModel)(nil)

type (
	// HomepageSlotSettingsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomepageSlotSettingsModel.
	HomepageSlotSettingsModel interface {
		homepageSlotSettingsModel
		withSession(session sqlx.Session) HomepageSlotSettingsModel
	}

	customHomepageSlotSettingsModel struct {
		*defaultHomepageSlotSettingsModel
	}
)

// NewHomepageSlotSettingsModel returns a model for the database table.
func NewHomepageSlotSettingsModel(conn sqlx.SqlConn) HomepageSlotSettingsModel {
	return &customHomepageSlotSettingsModel{
		defaultHomepageSlotSettingsModel: newHomepageSlotSettingsModel(conn),
	}
}

func (m *customHomepageSlotSettingsModel) withSession(session sqlx.Session) HomepageSlotSettingsModel {
	return NewHomepageSlotSettingsModel(sqlx.NewSqlConnFromSession(session))
}
