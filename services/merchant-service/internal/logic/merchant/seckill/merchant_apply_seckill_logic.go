package seckill

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantApplySeckillLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantApplySeckillLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantApplySeckillLogic {
	return &MerchantApplySeckillLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantApplySeckillLogic) MerchantApplySeckill(ctx context.Context, req *types.SeckillApplyReq) (resp *types.AnyResp, err error) {
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	entry, err := biz.NewMerchantLogic(l.svcCtx).ApplySeckill(shopID, userID, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: entry}, nil
}
