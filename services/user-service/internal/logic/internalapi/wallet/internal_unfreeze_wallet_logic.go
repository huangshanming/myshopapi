package wallet

import (
	"context"
	"mymall/pkg/appinput"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalUnfreezeWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalUnfreezeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalUnfreezeWalletLogic {
	return &InternalUnfreezeWalletLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalUnfreezeWalletLogic) InternalUnfreezeWallet(ctx context.Context, req *types.WalletOrderOpReq) error {
	_, err := hinternal.NewWalletHandler(l.svcCtx).Unfreeze(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return err
	}
	return nil
}
