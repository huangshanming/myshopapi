package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ArticleFavoritesModel = (*customArticleFavoritesModel)(nil)

type (
	// ArticleFavoritesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleFavoritesModel.
	ArticleFavoritesModel interface {
		articleFavoritesModel
		withSession(session sqlx.Session) ArticleFavoritesModel
	}

	customArticleFavoritesModel struct {
		*defaultArticleFavoritesModel
	}
)

// NewArticleFavoritesModel returns a model for the database table.
func NewArticleFavoritesModel(conn sqlx.SqlConn) ArticleFavoritesModel {
	return &customArticleFavoritesModel{
		defaultArticleFavoritesModel: newArticleFavoritesModel(conn),
	}
}

func (m *customArticleFavoritesModel) withSession(session sqlx.Session) ArticleFavoritesModel {
	return NewArticleFavoritesModel(sqlx.NewSqlConnFromSession(session))
}
