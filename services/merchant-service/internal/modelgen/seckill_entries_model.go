package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SeckillEntriesModel = (*customSeckillEntriesModel)(nil)

type (
	// SeckillEntriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSeckillEntriesModel.
	SeckillEntriesModel interface {
		seckillEntriesModel
		withSession(session sqlx.Session) SeckillEntriesModel
	}

	customSeckillEntriesModel struct {
		*defaultSeckillEntriesModel
	}
)

// NewSeckillEntriesModel returns a model for the database table.
func NewSeckillEntriesModel(conn sqlx.SqlConn) SeckillEntriesModel {
	return &customSeckillEntriesModel{
		defaultSeckillEntriesModel: newSeckillEntriesModel(conn),
	}
}

func (m *customSeckillEntriesModel) withSession(session sqlx.Session) SeckillEntriesModel {
	return NewSeckillEntriesModel(sqlx.NewSqlConnFromSession(session))
}
