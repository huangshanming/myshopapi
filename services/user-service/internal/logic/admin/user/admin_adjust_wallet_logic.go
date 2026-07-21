package user

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

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

func (l *AdminAdjustWalletLogic) AdminAdjustWallet(ctx context.Context, req *types.WalletAdjustReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/users/{Id}/wallet/adjust", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewWalletHandler(l.svcCtx).AdminAdjustWallet)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
