package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserGetWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetWalletLogic {
	return &UserGetWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetWalletLogic) UserGetWallet(w http.ResponseWriter, r *http.Request) {
	huser.NewWalletHandler(l.svcCtx).UserGetWallet(w, r)
}
