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

type AdminGetWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGetWalletLogic(svcCtx *svc.ServiceContext) *AdminGetWalletLogic {
	return &AdminGetWalletLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminGetWalletLogic) AdminGetWallet(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/shops/:id/wallet", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hadmin.NewWalletHandler(l.svcCtx).AdminGetWallet)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
