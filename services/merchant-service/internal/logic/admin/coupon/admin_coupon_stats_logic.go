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

type AdminCouponStatsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCouponStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCouponStatsLogic {
	return &AdminCouponStatsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCouponStatsLogic) AdminCouponStats(ctx context.Context, req *types.IdPathReq) (resp *types.CouponStatsResp, err error) {
	id := req.Id
	st, err := biz.NewMerchantLogic(l.svcCtx).CouponStats(id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.CouponStatsResp{Data: st}, nil
}
