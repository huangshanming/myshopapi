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

type MerchantAuthMeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantAuthMeLogic(svcCtx *svc.ServiceContext) *MerchantAuthMeLogic {
	return &MerchantAuthMeLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantAuthMeLogic) MerchantAuthMe(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/merchant/auth/me", nil, nil, nil, hhandler.NewShopOpsHandler(l.svcCtx).AuthMe)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
