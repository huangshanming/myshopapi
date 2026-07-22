package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CommunityArticleModel = (*customCommunityArticleModel)(nil)

type (
	// CommunityArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommunityArticleModel.
	CommunityArticleModel interface {
		communityArticleModel
		withSession(session sqlx.Session) CommunityArticleModel
	}

	customCommunityArticleModel struct {
		*defaultCommunityArticleModel
	}
)

// NewCommunityArticleModel returns a model for the database table.
func NewCommunityArticleModel(conn sqlx.SqlConn) CommunityArticleModel {
	return &customCommunityArticleModel{
		defaultCommunityArticleModel: newCommunityArticleModel(conn),
	}
}

func (m *customCommunityArticleModel) withSession(session sqlx.Session) CommunityArticleModel {
	return NewCommunityArticleModel(sqlx.NewSqlConnFromSession(session))
}
