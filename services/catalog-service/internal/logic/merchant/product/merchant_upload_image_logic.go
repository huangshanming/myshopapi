package product

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	hmerchant "mymall/services/catalog-service/internal/product/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantUploadImageLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantUploadImageLogic {
	return &MerchantUploadImageLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUploadImageLogic) MerchantUploadImage(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	data, err := hmerchant.NewProductHandler(l.svcCtx).Upload(ctx, appinput.CallInput{Request: r})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
