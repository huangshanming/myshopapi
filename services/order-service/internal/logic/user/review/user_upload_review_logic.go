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

type UserUploadReviewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserUploadReviewLogic(svcCtx *svc.ServiceContext) *UserUploadReviewLogic {
	return &UserUploadReviewLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserUploadReviewLogic) UserUploadReview(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/review-uploads", nil, nil, req, huser.NewReviewHandler(l.svcCtx).Upload)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
