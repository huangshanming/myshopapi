package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ArticleAudiencesModel = (*customArticleAudiencesModel)(nil)

type (
	// ArticleAudiencesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleAudiencesModel.
	ArticleAudiencesModel interface {
		articleAudiencesModel
		withSession(session sqlx.Session) ArticleAudiencesModel
	}

	customArticleAudiencesModel struct {
		*defaultArticleAudiencesModel
	}
)

// NewArticleAudiencesModel returns a model for the database table.
func NewArticleAudiencesModel(conn sqlx.SqlConn) ArticleAudiencesModel {
	return &customArticleAudiencesModel{
		defaultArticleAudiencesModel: newArticleAudiencesModel(conn),
	}
}

func (m *customArticleAudiencesModel) withSession(session sqlx.Session) ArticleAudiencesModel {
	return NewArticleAudiencesModel(sqlx.NewSqlConnFromSession(session))
}
