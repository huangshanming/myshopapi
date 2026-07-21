package seckill

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillConsumeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSeckillConsumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillConsumeLogic {
	return &SeckillConsumeLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *SeckillConsumeLogic) SeckillConsume(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body types.SeckillConsumeReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	data, err := biz.NewMerchantLogic(l.svcCtx).ConsumeSeckill(body.EntryID, body.ProductID, body.Quantity)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
