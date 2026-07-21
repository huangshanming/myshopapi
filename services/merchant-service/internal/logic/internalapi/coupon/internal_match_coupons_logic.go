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

type InternalMatchCouponsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalMatchCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalMatchCouponsLogic {
	return &InternalMatchCouponsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalMatchCouponsLogic) InternalMatchCoupons(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hinternal.NewCouponHandler(l.svcCtx).InternalMatchCoupons(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
