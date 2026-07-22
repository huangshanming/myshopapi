package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CouponRedeemLogsModel = (*customCouponRedeemLogsModel)(nil)

type (
	// CouponRedeemLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponRedeemLogsModel.
	CouponRedeemLogsModel interface {
		couponRedeemLogsModel
		withSession(session sqlx.Session) CouponRedeemLogsModel
	}

	customCouponRedeemLogsModel struct {
		*defaultCouponRedeemLogsModel
	}
)

// NewCouponRedeemLogsModel returns a model for the database table.
func NewCouponRedeemLogsModel(conn sqlx.SqlConn) CouponRedeemLogsModel {
	return &customCouponRedeemLogsModel{
		defaultCouponRedeemLogsModel: newCouponRedeemLogsModel(conn),
	}
}

func (m *customCouponRedeemLogsModel) withSession(session sqlx.Session) CouponRedeemLogsModel {
	return NewCouponRedeemLogsModel(sqlx.NewSqlConnFromSession(session))
}
