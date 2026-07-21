package wallet

import (
	"context"
	"mymall/pkg/httpinvoke"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalUnfreezeWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalUnfreezeWalletLogic(svcCtx *svc.ServiceContext) *InternalUnfreezeWalletLogic {
	return &InternalUnfreezeWalletLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *InternalUnfreezeWalletLogic) InternalUnfreezeWallet(ctx context.Context, req *types.WalletOrderOpReq) error {
	_, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/wallet/unfreeze", nil, nil, req, hinternal.NewWalletHandler(l.svcCtx).Unfreeze)
	if err != nil {
		return err
	}
	return nil
}
