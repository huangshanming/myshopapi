package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CouponGrantsModel = (*customCouponGrantsModel)(nil)

type (
	// CouponGrantsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponGrantsModel.
	CouponGrantsModel interface {
		couponGrantsModel
		withSession(session sqlx.Session) CouponGrantsModel
	}

	customCouponGrantsModel struct {
		*defaultCouponGrantsModel
	}
)

// NewCouponGrantsModel returns a model for the database table.
func NewCouponGrantsModel(conn sqlx.SqlConn) CouponGrantsModel {
	return &customCouponGrantsModel{
		defaultCouponGrantsModel: newCouponGrantsModel(conn),
	}
}

func (m *customCouponGrantsModel) withSession(session sqlx.Session) CouponGrantsModel {
	return NewCouponGrantsModel(sqlx.NewSqlConnFromSession(session))
}
