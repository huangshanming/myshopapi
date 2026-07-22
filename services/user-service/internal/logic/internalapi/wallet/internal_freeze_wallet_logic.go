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

type InternalFreezeWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalFreezeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalFreezeWalletLogic {
	return &InternalFreezeWalletLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalFreezeWalletLogic) InternalFreezeWallet(ctx context.Context, req *types.WalletOrderOpReq) (*types.EmptyResp, error) {
	if err := biz.NewWalletLogic(l.svcCtx).FreezeForOrder(ctx, req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
