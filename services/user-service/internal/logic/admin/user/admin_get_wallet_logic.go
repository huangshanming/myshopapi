package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type AdminGetWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetWalletLogic {
	return &AdminGetWalletLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGetWalletLogic) AdminGetWallet(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	wallet, err := biz.NewWalletLogic(l.svcCtx).GetWallet(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: wallet}, nil
}
