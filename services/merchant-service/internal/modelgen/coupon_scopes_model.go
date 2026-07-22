package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CouponScopesModel = (*customCouponScopesModel)(nil)

type (
	// CouponScopesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponScopesModel.
	CouponScopesModel interface {
		couponScopesModel
		withSession(session sqlx.Session) CouponScopesModel
	}

	customCouponScopesModel struct {
		*defaultCouponScopesModel
	}
)

// NewCouponScopesModel returns a model for the database table.
func NewCouponScopesModel(conn sqlx.SqlConn) CouponScopesModel {
	return &customCouponScopesModel{
		defaultCouponScopesModel: newCouponScopesModel(conn),
	}
}

func (m *customCouponScopesModel) withSession(session sqlx.Session) CouponScopesModel {
	return NewCouponScopesModel(sqlx.NewSqlConnFromSession(session))
}
