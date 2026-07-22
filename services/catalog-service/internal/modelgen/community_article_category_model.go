package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CommunityArticleCategoryModel = (*customCommunityArticleCategoryModel)(nil)

type (
	// CommunityArticleCategoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommunityArticleCategoryModel.
	CommunityArticleCategoryModel interface {
		communityArticleCategoryModel
		withSession(session sqlx.Session) CommunityArticleCategoryModel
	}

	customCommunityArticleCategoryModel struct {
		*defaultCommunityArticleCategoryModel
	}
)

// NewCommunityArticleCategoryModel returns a model for the database table.
func NewCommunityArticleCategoryModel(conn sqlx.SqlConn) CommunityArticleCategoryModel {
	return &customCommunityArticleCategoryModel{
		defaultCommunityArticleCategoryModel: newCommunityArticleCategoryModel(conn),
	}
}

func (m *customCommunityArticleCategoryModel) withSession(session sqlx.Session) CommunityArticleCategoryModel {
	return NewCommunityArticleCategoryModel(sqlx.NewSqlConnFromSession(session))
}
