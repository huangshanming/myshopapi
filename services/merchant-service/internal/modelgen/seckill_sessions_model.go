package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SeckillSessionsModel = (*customSeckillSessionsModel)(nil)

type (
	// SeckillSessionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSeckillSessionsModel.
	SeckillSessionsModel interface {
		seckillSessionsModel
		withSession(session sqlx.Session) SeckillSessionsModel
	}

	customSeckillSessionsModel struct {
		*defaultSeckillSessionsModel
	}
)

// NewSeckillSessionsModel returns a model for the database table.
func NewSeckillSessionsModel(conn sqlx.SqlConn) SeckillSessionsModel {
	return &customSeckillSessionsModel{
		defaultSeckillSessionsModel: newSeckillSessionsModel(conn),
	}
}

func (m *customSeckillSessionsModel) withSession(session sqlx.Session) SeckillSessionsModel {
	return NewSeckillSessionsModel(sqlx.NewSqlConnFromSession(session))
}
