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

type AdminUpdateSeckillRuleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateSeckillRuleLogic(svcCtx *svc.ServiceContext) *AdminUpdateSeckillRuleLogic {
	return &AdminUpdateSeckillRuleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateSeckillRuleLogic) AdminUpdateSeckillRule(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/seckill/rule", nil, nil, req, hadmin.NewSeckillHandler(l.svcCtx).AdminUpdateSeckillRule)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
