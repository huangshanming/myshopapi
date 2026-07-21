package seckill

import (
	"context"
	"mymall/pkg/appinput"
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

func (l *MerchantApplySeckillLogic) MerchantApplySeckill(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	var body types.SeckillApplyReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	entry, err := biz.NewMerchantLogic(l.svcCtx).ApplySeckill(shopID, userID, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: entry}, nil
}
