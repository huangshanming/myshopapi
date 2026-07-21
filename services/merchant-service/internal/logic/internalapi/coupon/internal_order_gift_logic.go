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

type InternalOrderGiftLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalOrderGiftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalOrderGiftLogic {
	return &InternalOrderGiftLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalOrderGiftLogic) InternalOrderGift(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hinternal.NewCouponHandler(l.svcCtx).InternalOrderGift(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
