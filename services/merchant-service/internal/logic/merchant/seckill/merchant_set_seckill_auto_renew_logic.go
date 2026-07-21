package seckill

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantSetSeckillAutoRenewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantSetSeckillAutoRenewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantSetSeckillAutoRenewLogic {
	return &MerchantSetSeckillAutoRenewLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantSetSeckillAutoRenewLogic) MerchantSetSeckillAutoRenew(ctx context.Context, req *types.SeckillAutoRenewBodyReq) (resp *types.AnyResp, err error) {
	shopID := middleware.GetShopID(ctx)
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "报名ID无效")
	}
	entry, err := biz.NewMerchantLogic(l.svcCtx).SetSeckillAutoRenew(shopID, req.Id, req.AutoRenew)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: entry}, nil
}
