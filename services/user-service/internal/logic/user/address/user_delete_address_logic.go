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

type UserDeleteAddressLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteAddressLogic {
	return &UserDeleteAddressLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteAddressLogic) UserDeleteAddress(ctx context.Context, req *types.IdPathReq) error {
	_, err := huser.NewAddressHandler(l.svcCtx).Delete(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}})
	if err != nil {
		return err
	}
	return nil
}
