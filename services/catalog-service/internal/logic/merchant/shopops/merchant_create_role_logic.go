package shopops

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hhandler "mymall/services/catalog-service/internal/shopops/handler"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantCreateRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCreateRoleLogic {
	return &MerchantCreateRoleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateRoleLogic) MerchantCreateRole(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hhandler.NewShopOpsHandler(l.svcCtx).SaveRole(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
