package coupon

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hpublic "mymall/services/merchant-service/internal/app/public"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicCouponPopupLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicCouponPopupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicCouponPopupLogic {
	return &PublicCouponPopupLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicCouponPopupLogic) PublicCouponPopup(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hpublic.NewCouponHandler(l.svcCtx).PublicCouponPopup(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
