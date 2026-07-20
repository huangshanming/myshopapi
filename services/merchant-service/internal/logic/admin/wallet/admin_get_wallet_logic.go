package wallet

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"
)

type AdminGetWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetWalletLogic {
	return &AdminGetWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetWalletLogic) AdminGetWallet(w http.ResponseWriter, r *http.Request) {
	hadmin.NewWalletHandler(l.svcCtx).AdminGetWallet(w, r)
}
