package wallet

import (
	"context"
	"mymall/pkg/httpinvoke"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalFreezeWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalFreezeWalletLogic(svcCtx *svc.ServiceContext) *InternalFreezeWalletLogic {
	return &InternalFreezeWalletLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *InternalFreezeWalletLogic) InternalFreezeWallet(ctx context.Context, req *types.WalletOrderOpReq) error {
	_, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/wallet/freeze", nil, nil, req, hinternal.NewWalletHandler(l.svcCtx).Freeze)
	if err != nil {
		return err
	}
	return nil
}
