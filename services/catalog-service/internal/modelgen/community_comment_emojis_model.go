package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CommunityCommentEmojisModel = (*customCommunityCommentEmojisModel)(nil)

type (
	// CommunityCommentEmojisModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommunityCommentEmojisModel.
	CommunityCommentEmojisModel interface {
		communityCommentEmojisModel
		withSession(session sqlx.Session) CommunityCommentEmojisModel
	}

	customCommunityCommentEmojisModel struct {
		*defaultCommunityCommentEmojisModel
	}
)

// NewCommunityCommentEmojisModel returns a model for the database table.
func NewCommunityCommentEmojisModel(conn sqlx.SqlConn) CommunityCommentEmojisModel {
	return &customCommunityCommentEmojisModel{
		defaultCommunityCommentEmojisModel: newCommunityCommentEmojisModel(conn),
	}
}

func (m *customCommunityCommentEmojisModel) withSession(session sqlx.Session) CommunityCommentEmojisModel {
	return NewCommunityCommentEmojisModel(sqlx.NewSqlConnFromSession(session))
}
