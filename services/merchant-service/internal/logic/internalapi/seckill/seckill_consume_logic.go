package seckill

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hinternal "mymall/services/merchant-service/internal/app/internalapi"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillConsumeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSeckillConsumeLogic(svcCtx *svc.ServiceContext) *SeckillConsumeLogic {
	return &SeckillConsumeLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *SeckillConsumeLogic) SeckillConsume(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/seckill/consume", nil, nil, req, hinternal.NewSeckillHandler(l.svcCtx).SeckillConsume)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
