package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CommunityArticleImgModel = (*customCommunityArticleImgModel)(nil)

type (
	// CommunityArticleImgModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommunityArticleImgModel.
	CommunityArticleImgModel interface {
		communityArticleImgModel
		withSession(session sqlx.Session) CommunityArticleImgModel
	}

	customCommunityArticleImgModel struct {
		*defaultCommunityArticleImgModel
	}
)

// NewCommunityArticleImgModel returns a model for the database table.
func NewCommunityArticleImgModel(conn sqlx.SqlConn) CommunityArticleImgModel {
	return &customCommunityArticleImgModel{
		defaultCommunityArticleImgModel: newCommunityArticleImgModel(conn),
	}
}

func (m *customCommunityArticleImgModel) withSession(session sqlx.Session) CommunityArticleImgModel {
	return NewCommunityArticleImgModel(sqlx.NewSqlConnFromSession(session))
}
