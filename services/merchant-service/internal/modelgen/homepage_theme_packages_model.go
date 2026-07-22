package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomepageThemePackagesModel = (*customHomepageThemePackagesModel)(nil)

type (
	// HomepageThemePackagesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomepageThemePackagesModel.
	HomepageThemePackagesModel interface {
		homepageThemePackagesModel
		withSession(session sqlx.Session) HomepageThemePackagesModel
	}

	customHomepageThemePackagesModel struct {
		*defaultHomepageThemePackagesModel
	}
)

// NewHomepageThemePackagesModel returns a model for the database table.
func NewHomepageThemePackagesModel(conn sqlx.SqlConn) HomepageThemePackagesModel {
	return &customHomepageThemePackagesModel{
		defaultHomepageThemePackagesModel: newHomepageThemePackagesModel(conn),
	}
}

func (m *customHomepageThemePackagesModel) withSession(session sqlx.Session) HomepageThemePackagesModel {
	return NewHomepageThemePackagesModel(sqlx.NewSqlConnFromSession(session))
}
