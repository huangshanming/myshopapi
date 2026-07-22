package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SeckillRulesModel = (*customSeckillRulesModel)(nil)

type (
	// SeckillRulesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSeckillRulesModel.
	SeckillRulesModel interface {
		seckillRulesModel
		withSession(session sqlx.Session) SeckillRulesModel
	}

	customSeckillRulesModel struct {
		*defaultSeckillRulesModel
	}
)

// NewSeckillRulesModel returns a model for the database table.
func NewSeckillRulesModel(conn sqlx.SqlConn) SeckillRulesModel {
	return &customSeckillRulesModel{
		defaultSeckillRulesModel: newSeckillRulesModel(conn),
	}
}

func (m *customSeckillRulesModel) withSession(session sqlx.Session) SeckillRulesModel {
	return NewSeckillRulesModel(sqlx.NewSqlConnFromSession(session))
}
