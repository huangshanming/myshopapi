package review

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	huser "mymall/services/order-service/internal/app/user"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicListProductReviewsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicListProductReviewsLogic(svcCtx *svc.ServiceContext) *PublicListProductReviewsLogic {
	return &PublicListProductReviewsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *PublicListProductReviewsLogic) PublicListProductReviews(ctx context.Context, req *types.IdPathReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/products/:id/reviews", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, huser.NewReviewHandler(l.svcCtx).ProductList)
	if err != nil {
		return nil, err
	}
	var out types.PageListResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		var list interface{}
		if err2 := httpinvoke.Decode(raw, &list); err2 == nil {
			return &types.PageListResp{List: list}, nil
		}
		return nil, err
	}
	return &out, nil
}
