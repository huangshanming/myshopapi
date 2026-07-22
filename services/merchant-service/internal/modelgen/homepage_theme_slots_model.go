package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomepageThemeSlotsModel = (*customHomepageThemeSlotsModel)(nil)

type (
	// HomepageThemeSlotsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomepageThemeSlotsModel.
	HomepageThemeSlotsModel interface {
		homepageThemeSlotsModel
		withSession(session sqlx.Session) HomepageThemeSlotsModel
	}

	customHomepageThemeSlotsModel struct {
		*defaultHomepageThemeSlotsModel
	}
)

// NewHomepageThemeSlotsModel returns a model for the database table.
func NewHomepageThemeSlotsModel(conn sqlx.SqlConn) HomepageThemeSlotsModel {
	return &customHomepageThemeSlotsModel{
		defaultHomepageThemeSlotsModel: newHomepageThemeSlotsModel(conn),
	}
}

func (m *customHomepageThemeSlotsModel) withSession(session sqlx.Session) HomepageThemeSlotsModel {
	return NewHomepageThemeSlotsModel(sqlx.NewSqlConnFromSession(session))
}
