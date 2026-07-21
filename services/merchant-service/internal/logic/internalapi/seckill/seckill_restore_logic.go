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

type SeckillRestoreLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSeckillRestoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillRestoreLogic {
	return &SeckillRestoreLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *SeckillRestoreLogic) SeckillRestore(ctx context.Context, req *types.SeckillRestoreReq) (resp *types.AnyResp, err error) {
	if err := biz.NewMerchantLogic(l.svcCtx).RestoreSeckill(req.EntryID, req.Quantity); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
