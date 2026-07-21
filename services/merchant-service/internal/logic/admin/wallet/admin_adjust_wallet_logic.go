package wallet

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

type AdminAdjustWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminAdjustWalletLogic(svcCtx *svc.ServiceContext) *AdminAdjustWalletLogic {
	return &AdminAdjustWalletLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminAdjustWalletLogic) AdminAdjustWallet(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/shops/:id/wallet/adjust", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewWalletHandler(l.svcCtx).AdminAdjustWallet)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
