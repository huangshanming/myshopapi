package user

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminAdjustWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminAdjustWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAdjustWalletLogic {
	return &AdminAdjustWalletLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminAdjustWalletLogic) AdminAdjustWallet(ctx context.Context, req *types.WalletAdjustReq) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewWalletHandler(l.svcCtx).AdminAdjustWallet(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
