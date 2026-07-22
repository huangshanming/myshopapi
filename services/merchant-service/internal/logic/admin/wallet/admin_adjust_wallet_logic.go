package wallet

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *AdminAdjustWalletLogic) AdminAdjustWallet(ctx context.Context, req *types.WalletAdjustBodyReq) (resp *types.WalletResp, err error) {
	shopID := req.Id
	adminID, _ := middleware.GetUserID(ctx)
	wallet, err := biz.NewMerchantLogic(l.svcCtx).AdjustWallet(shopID, req.Field, req.Amount, req.Remark, adminID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.WalletResp{Data: wallet}, nil
}
