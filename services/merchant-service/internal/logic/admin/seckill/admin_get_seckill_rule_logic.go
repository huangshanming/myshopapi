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

func (l *AdminGetSeckillRuleLogic) AdminGetSeckillRule(ctx context.Context) (resp *types.SeckillRuleResp, err error) {

	rule, err := biz.NewMerchantLogic(l.svcCtx).GetSeckillRule()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.SeckillRuleResp{Data: rule}, nil
}
