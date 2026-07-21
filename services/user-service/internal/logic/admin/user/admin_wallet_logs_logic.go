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

type AdminWalletLogsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminWalletLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminWalletLogsLogic {
	return &AdminWalletLogsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminWalletLogsLogic) AdminWalletLogs(ctx context.Context, req *types.AdminWalletLogsReq) (resp *types.PageListResp, err error) {
	list, total, err := biz.NewWalletLogic(l.svcCtx).ListWalletLogs(ctx, req.Id, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
