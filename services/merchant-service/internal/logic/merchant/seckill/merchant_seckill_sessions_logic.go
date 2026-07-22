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

type MerchantSeckillSessionsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantSeckillSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantSeckillSessionsLogic {
	return &MerchantSeckillSessionsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantSeckillSessionsLogic) MerchantSeckillSessions(ctx context.Context) (resp *types.SeckillSessionsResp, err error) {

	data, err := biz.NewMerchantLogic(l.svcCtx).MerchantSeckillSessions()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.SeckillSessionsResp{Data: data}, nil
}
