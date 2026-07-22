package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ArticleLikesModel = (*customArticleLikesModel)(nil)

type (
	// ArticleLikesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleLikesModel.
	ArticleLikesModel interface {
		articleLikesModel
		withSession(session sqlx.Session) ArticleLikesModel
	}

	customArticleLikesModel struct {
		*defaultArticleLikesModel
	}
)

// NewArticleLikesModel returns a model for the database table.
func NewArticleLikesModel(conn sqlx.SqlConn) ArticleLikesModel {
	return &customArticleLikesModel{
		defaultArticleLikesModel: newArticleLikesModel(conn),
	}
}

func (m *customArticleLikesModel) withSession(session sqlx.Session) ArticleLikesModel {
	return NewArticleLikesModel(sqlx.NewSqlConnFromSession(session))
}
