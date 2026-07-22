package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomepageBannersModel = (*customHomepageBannersModel)(nil)

type (
	// HomepageBannersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomepageBannersModel.
	HomepageBannersModel interface {
		homepageBannersModel
		withSession(session sqlx.Session) HomepageBannersModel
	}

	customHomepageBannersModel struct {
		*defaultHomepageBannersModel
	}
)

// NewHomepageBannersModel returns a model for the database table.
func NewHomepageBannersModel(conn sqlx.SqlConn) HomepageBannersModel {
	return &customHomepageBannersModel{
		defaultHomepageBannersModel: newHomepageBannersModel(conn),
	}
}

func (m *customHomepageBannersModel) withSession(session sqlx.Session) HomepageBannersModel {
	return NewHomepageBannersModel(sqlx.NewSqlConnFromSession(session))
}
