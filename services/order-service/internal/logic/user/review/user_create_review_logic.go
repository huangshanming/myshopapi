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

type UserCreateReviewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCreateReviewLogic(svcCtx *svc.ServiceContext) *UserCreateReviewLogic {
	return &UserCreateReviewLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserCreateReviewLogic) UserCreateReview(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/orders/:id/reviews", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, huser.NewReviewHandler(l.svcCtx).Create)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
