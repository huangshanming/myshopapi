package address

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetDefaultLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSetDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultLogic {
	return &SetDefaultLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *SetDefaultLogic) SetDefault(ctx context.Context, req *types.IdPathReq) (*types.EmptyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	if err := biz.NewAddressLogic(l.svcCtx).SetDefault(ctx, userID, req.Id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
