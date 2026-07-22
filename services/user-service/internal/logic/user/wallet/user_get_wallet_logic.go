package wallet

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

type UserGetWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetWalletLogic {
	return &UserGetWalletLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserGetWalletLogic) UserGetWallet(ctx context.Context) (*types.WalletResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	w, err := biz.NewWalletLogic(l.svcCtx).GetWallet(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.WalletResp{Data: w}, nil
}
