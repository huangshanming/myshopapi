package seckill

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

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

func (l *MerchantSetSeckillAutoRenewLogic) MerchantSetSeckillAutoRenew(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	shopID := middleware.GetShopID(ctx)
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "报名ID无效")
	}
	var body types.SeckillAutoRenewReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	entry, err := biz.NewMerchantLogic(l.svcCtx).SetSeckillAutoRenew(shopID, id, body.AutoRenew)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: entry}, nil
}
