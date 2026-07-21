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

type MerchantAuthMeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantAuthMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantAuthMeLogic {
	return &MerchantAuthMeLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantAuthMeLogic) MerchantAuthMe(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hhandler.NewShopOpsHandler(l.svcCtx).AuthMe(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
