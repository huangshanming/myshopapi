package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductReviewsModel = (*customProductReviewsModel)(nil)

type (
	// ProductReviewsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductReviewsModel.
	ProductReviewsModel interface {
		productReviewsModel
		withSession(session sqlx.Session) ProductReviewsModel
	}

	customProductReviewsModel struct {
		*defaultProductReviewsModel
	}
)

// NewProductReviewsModel returns a model for the database table.
func NewProductReviewsModel(conn sqlx.SqlConn) ProductReviewsModel {
	return &customProductReviewsModel{
		defaultProductReviewsModel: newProductReviewsModel(conn),
	}
}

func (m *customProductReviewsModel) withSession(session sqlx.Session) ProductReviewsModel {
	return NewProductReviewsModel(sqlx.NewSqlConnFromSession(session))
}
