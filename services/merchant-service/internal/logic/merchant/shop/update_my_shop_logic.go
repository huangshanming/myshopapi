package shop

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hmerchant "mymall/services/merchant-service/internal/app/merchant"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMyShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateMyShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyShopLogic {
	return &UpdateMyShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyShopLogic) UpdateMyShop(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hmerchant.NewShopHandler(l.svcCtx).UpdateMyShop(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
