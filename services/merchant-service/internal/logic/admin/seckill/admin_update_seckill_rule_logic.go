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

type AdminUpdateSeckillRuleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateSeckillRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateSeckillRuleLogic {
	return &AdminUpdateSeckillRuleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateSeckillRuleLogic) AdminUpdateSeckillRule(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body types.SeckillRuleReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	rule, err := biz.NewMerchantLogic(l.svcCtx).UpdateSeckillRule(body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: rule}, nil
}
