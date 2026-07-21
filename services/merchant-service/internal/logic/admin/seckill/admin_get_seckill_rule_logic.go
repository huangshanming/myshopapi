package seckill

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hadmin "mymall/services/merchant-service/internal/app/admin"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetSeckillRuleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGetSeckillRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetSeckillRuleLogic {
	return &AdminGetSeckillRuleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGetSeckillRuleLogic) AdminGetSeckillRule(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hadmin.NewSeckillHandler(l.svcCtx).AdminGetSeckillRule(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
