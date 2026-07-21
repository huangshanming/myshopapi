package review

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	huser "mymall/services/order-service/internal/app/user"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUploadReviewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserUploadReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUploadReviewLogic {
	return &UserUploadReviewLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserUploadReviewLogic) UserUploadReview(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	data, err := huser.NewReviewHandler(l.svcCtx).Upload(ctx, appinput.CallInput{Request: r})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
