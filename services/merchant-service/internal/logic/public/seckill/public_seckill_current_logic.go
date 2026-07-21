package seckill

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hpublic "mymall/services/merchant-service/internal/app/public"
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
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hpublic.NewSeckillHandler(l.svcCtx).PublicSeckillCurrent(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
