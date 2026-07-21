package address

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateAddressLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateAddressLogic {
	return &UserUpdateAddressLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateAddressLogic) UserUpdateAddress(ctx context.Context, req *types.AddressUpdateReq) error {
	_, err := huser.NewAddressHandler(l.svcCtx).Update(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}, Body: req})
	if err != nil {
		return err
	}
	return nil
}
