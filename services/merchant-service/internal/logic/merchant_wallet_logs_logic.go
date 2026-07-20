package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantWalletLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantWalletLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantWalletLogsLogic {
	return &MerchantWalletLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantWalletLogsLogic) MerchantWalletLogs(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewWalletHandler(l.svcCtx).MerchantWalletLogs(w, r)
}
