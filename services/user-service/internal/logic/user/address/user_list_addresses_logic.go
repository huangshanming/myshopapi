package address

import (
	"context"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListAddressesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListAddressesLogic {
	return &UserListAddressesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserListAddressesLogic) UserListAddresses(ctx context.Context) (resp *types.PageListResp, err error) {
	data, err := huser.NewAddressHandler(l.svcCtx).List(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.PageListResp{List: data}, nil
}
