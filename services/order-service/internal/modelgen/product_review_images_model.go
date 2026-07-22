package modelgen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductReviewImagesModel = (*customProductReviewImagesModel)(nil)

type (
	// ProductReviewImagesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductReviewImagesModel.
	ProductReviewImagesModel interface {
		productReviewImagesModel
		withSession(session sqlx.Session) ProductReviewImagesModel
	}

	customProductReviewImagesModel struct {
		*defaultProductReviewImagesModel
	}
)

// NewProductReviewImagesModel returns a model for the database table.
func NewProductReviewImagesModel(conn sqlx.SqlConn) ProductReviewImagesModel {
	return &customProductReviewImagesModel{
		defaultProductReviewImagesModel: newProductReviewImagesModel(conn),
	}
}

func (m *customProductReviewImagesModel) withSession(session sqlx.Session) ProductReviewImagesModel {
	return NewProductReviewImagesModel(sqlx.NewSqlConnFromSession(session))
}
