package review

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicListProductReviewsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicListProductReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicListProductReviewsLogic {
	return &PublicListProductReviewsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *PublicListProductReviewsLogic) PublicListProductReviews(ctx context.Context, req *types.IdPathReq) (*types.PageListResp, error) {
	list, total, err := biz.NewReviewLogic(l.svcCtx).ListByProduct(ctx, req.Id, 1, 10)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
