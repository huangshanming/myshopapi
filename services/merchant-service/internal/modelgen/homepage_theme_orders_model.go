package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomepageThemeOrdersModel = (*customHomepageThemeOrdersModel)(nil)

type (
	// HomepageThemeOrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomepageThemeOrdersModel.
	HomepageThemeOrdersModel interface {
		homepageThemeOrdersModel
		withSession(session sqlx.Session) HomepageThemeOrdersModel
	}

	customHomepageThemeOrdersModel struct {
		*defaultHomepageThemeOrdersModel
	}
)

// NewHomepageThemeOrdersModel returns a model for the database table.
func NewHomepageThemeOrdersModel(conn sqlx.SqlConn) HomepageThemeOrdersModel {
	return &customHomepageThemeOrdersModel{
		defaultHomepageThemeOrdersModel: newHomepageThemeOrdersModel(conn),
	}
}

func (m *customHomepageThemeOrdersModel) withSession(session sqlx.Session) HomepageThemeOrdersModel {
	return NewHomepageThemeOrdersModel(sqlx.NewSqlConnFromSession(session))
}
