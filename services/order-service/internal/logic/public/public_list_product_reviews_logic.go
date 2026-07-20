package public

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type PublicListProductReviewsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicListProductReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicListProductReviewsLogic {
	return &PublicListProductReviewsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicListProductReviewsLogic) PublicListProductReviews(w http.ResponseWriter, r *http.Request) {
	huser.NewReviewHandler(l.svcCtx).ProductList(w, r)
}
