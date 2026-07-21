package wallet

import (
	"context"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetWalletLogic {
	return &UserGetWalletLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserGetWalletLogic) UserGetWallet(ctx context.Context) (resp *types.AnyResp, err error) {
	data, err := huser.NewWalletHandler(l.svcCtx).UserGetWallet(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
