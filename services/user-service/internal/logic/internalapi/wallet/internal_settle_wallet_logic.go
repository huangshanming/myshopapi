package wallet

import (
	"context"
	"mymall/pkg/httpinvoke"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalSettleWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalSettleWalletLogic(svcCtx *svc.ServiceContext) *InternalSettleWalletLogic {
	return &InternalSettleWalletLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *InternalSettleWalletLogic) InternalSettleWallet(ctx context.Context, req *types.WalletOrderOpReq) error {
	_, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/wallet/settle", nil, nil, req, hinternal.NewWalletHandler(l.svcCtx).Settle)
	if err != nil {
		return err
	}
	return nil
}
