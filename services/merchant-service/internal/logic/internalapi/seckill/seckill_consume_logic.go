package seckill

import (
	"context"
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

func (l *SeckillConsumeLogic) SeckillConsume(ctx context.Context, req *types.SeckillConsumeReq) (resp *types.AnyResp, err error) {
	data, err := biz.NewMerchantLogic(l.svcCtx).ConsumeSeckill(req.EntryID, req.ProductID, req.Quantity)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
