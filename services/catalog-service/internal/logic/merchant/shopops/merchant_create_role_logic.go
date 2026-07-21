package shopops

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hhandler "mymall/services/catalog-service/internal/shopops/handler"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantCreateRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateRoleLogic(svcCtx *svc.ServiceContext) *MerchantCreateRoleLogic {
	return &MerchantCreateRoleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateRoleLogic) MerchantCreateRole(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/shop/roles", nil, nil, req, hhandler.NewShopOpsHandler(l.svcCtx).SaveRole)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
