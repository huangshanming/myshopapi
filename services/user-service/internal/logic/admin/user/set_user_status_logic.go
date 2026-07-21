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

type SetUserStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSetUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserStatusLogic {
	return &SetUserStatusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *SetUserStatusLogic) SetUserStatus(ctx context.Context, req *types.UserStatusReq) error {
	if err := biz.NewRBACLogic(l.svcCtx).SetUserStatus(ctx, req.Id, req.Status); err != nil {
		return xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil
}
