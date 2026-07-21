package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
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
	adminID, _ := middleware.GetUserID(ctx)
	wallet, err := biz.NewWalletLogic(l.svcCtx).AdjustWallet(ctx, req.Id, req.Field, req.Amount, req.Remark, adminID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: wallet}, nil
}
