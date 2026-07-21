package wallet

import (
	"context"
	"mymall/pkg/appinput"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalSettleWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalSettleWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalSettleWalletLogic {
	return &InternalSettleWalletLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalSettleWalletLogic) InternalSettleWallet(ctx context.Context, req *types.WalletOrderOpReq) error {
	_, err := hinternal.NewWalletHandler(l.svcCtx).Settle(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return err
	}
	return nil
}
