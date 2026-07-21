package seckill

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hadmin "mymall/services/merchant-service/internal/app/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetSeckillRuleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGetSeckillRuleLogic(svcCtx *svc.ServiceContext) *AdminGetSeckillRuleLogic {
	return &AdminGetSeckillRuleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminGetSeckillRuleLogic) AdminGetSeckillRule(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/seckill/rule", nil, nil, nil, hadmin.NewSeckillHandler(l.svcCtx).AdminGetSeckillRule)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
