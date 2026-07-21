package shop

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hmerchant "mymall/services/merchant-service/internal/app/merchant"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewApplyLogic(svcCtx *svc.ServiceContext) *ApplyLogic {
	return &ApplyLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ApplyLogic) Apply(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/apply", nil, nil, req, hmerchant.NewShopHandler(l.svcCtx).Apply)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
