package coupon

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hinternal "mymall/services/merchant-service/internal/app/internalapi"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalLockCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalLockCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalLockCouponLogic {
	return &InternalLockCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalLockCouponLogic) InternalLockCoupon(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hinternal.NewCouponHandler(l.svcCtx).InternalLockCoupon(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
