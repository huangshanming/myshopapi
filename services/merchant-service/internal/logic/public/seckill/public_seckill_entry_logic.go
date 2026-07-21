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

type PublicSeckillEntryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicSeckillEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicSeckillEntryLogic {
	return &PublicSeckillEntryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicSeckillEntryLogic) PublicSeckillEntry(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	data, err := biz.NewMerchantLogic(l.svcCtx).PublicSeckillEntry(id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
