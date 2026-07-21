package wallet

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserGetWalletLogic(svcCtx *svc.ServiceContext) *UserGetWalletLogic {
	return &UserGetWalletLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserGetWalletLogic) UserGetWallet(ctx context.Context) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/user/wallet", nil, nil, nil, huser.NewWalletHandler(l.svcCtx).UserGetWallet)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
