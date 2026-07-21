package address

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateAddressLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCreateAddressLogic(svcCtx *svc.ServiceContext) *UserCreateAddressLogic {
	return &UserCreateAddressLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserCreateAddressLogic) UserCreateAddress(ctx context.Context, req *types.AddressReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/addresses", nil, nil, req, huser.NewAddressHandler(l.svcCtx).Create)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
