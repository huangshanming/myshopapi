package seckill

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hpublic "mymall/services/merchant-service/internal/app/public"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicSeckillEntryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicSeckillEntryLogic(svcCtx *svc.ServiceContext) *PublicSeckillEntryLogic {
	return &PublicSeckillEntryLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *PublicSeckillEntryLogic) PublicSeckillEntry(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/seckill/entries/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hpublic.NewSeckillHandler(l.svcCtx).PublicSeckillEntry)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
