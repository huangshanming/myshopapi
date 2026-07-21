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

type PublicSeckillCurrentLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicSeckillCurrentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicSeckillCurrentLogic {
	return &PublicSeckillCurrentLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicSeckillCurrentLogic) PublicSeckillCurrent(ctx context.Context) (resp *types.AnyResp, err error) {

	data, err := biz.NewMerchantLogic(l.svcCtx).PublicSeckillCurrent()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
