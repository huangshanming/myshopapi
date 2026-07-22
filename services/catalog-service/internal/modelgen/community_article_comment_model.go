package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CommunityArticleCommentModel = (*customCommunityArticleCommentModel)(nil)

type (
	// CommunityArticleCommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommunityArticleCommentModel.
	CommunityArticleCommentModel interface {
		communityArticleCommentModel
		withSession(session sqlx.Session) CommunityArticleCommentModel
	}

	customCommunityArticleCommentModel struct {
		*defaultCommunityArticleCommentModel
	}
)

// NewCommunityArticleCommentModel returns a model for the database table.
func NewCommunityArticleCommentModel(conn sqlx.SqlConn) CommunityArticleCommentModel {
	return &customCommunityArticleCommentModel{
		defaultCommunityArticleCommentModel: newCommunityArticleCommentModel(conn),
	}
}

func (m *customCommunityArticleCommentModel) withSession(session sqlx.Session) CommunityArticleCommentModel {
	return NewCommunityArticleCommentModel(sqlx.NewSqlConnFromSession(session))
}
