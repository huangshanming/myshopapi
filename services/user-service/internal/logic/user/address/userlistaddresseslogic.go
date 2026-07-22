// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListAddressesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListAddressesLogic {
	return &UserListAddressesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListAddressesLogic) UserListAddresses() (resp *types.PageListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
