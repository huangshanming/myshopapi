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

func (l *SeckillRestoreLogic) SeckillRestore(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body types.SeckillRestoreReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := biz.NewMerchantLogic(l.svcCtx).RestoreSeckill(body.EntryID, body.Quantity); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
