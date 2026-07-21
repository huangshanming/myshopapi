package address

import (
	"context"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateAddressLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCreateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateAddressLogic {
	return &UserCreateAddressLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserCreateAddressLogic) UserCreateAddress(ctx context.Context, req *types.AddressReq) (resp *types.AnyResp, err error) {
	data, err := huser.NewAddressHandler(l.svcCtx).Create(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
