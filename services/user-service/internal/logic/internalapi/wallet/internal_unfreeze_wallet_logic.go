package wallet

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
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
	if err := biz.NewWalletLogic(l.svcCtx).UnfreezeOrder(ctx, req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		return xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil
}
