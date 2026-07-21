package address

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListAddressesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserListAddressesLogic(svcCtx *svc.ServiceContext) *UserListAddressesLogic {
	return &UserListAddressesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserListAddressesLogic) UserListAddresses(ctx context.Context) (resp *types.PageListResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/user/addresses", nil, nil, nil, huser.NewAddressHandler(l.svcCtx).List)
	if err != nil {
		return nil, err
	}
	var list interface{}
	if err := httpinvoke.Decode(raw, &list); err != nil {
		return nil, err
	}
	return &types.PageListResp{List: list}, nil
}
