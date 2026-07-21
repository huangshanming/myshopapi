package address

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateAddressLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserUpdateAddressLogic(svcCtx *svc.ServiceContext) *UserUpdateAddressLogic {
	return &UserUpdateAddressLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateAddressLogic) UserUpdateAddress(ctx context.Context, req *types.AddressUpdateReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/user/addresses/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, huser.NewAddressHandler(l.svcCtx).Update)
	if err != nil {
		return err
	}
	return nil
}
