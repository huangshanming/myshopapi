package coupon

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantCouponStatsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCouponStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCouponStatsLogic {
	return &MerchantCouponStatsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCouponStatsLogic) MerchantCouponStats(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	st, err := biz.NewMerchantLogic(l.svcCtx).CouponStats(id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: st}, nil
}
