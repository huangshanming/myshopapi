// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeleteAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteAddressLogic {
	return &UserDeleteAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteAddressLogic) UserDeleteAddress(req *types.IdPathReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
