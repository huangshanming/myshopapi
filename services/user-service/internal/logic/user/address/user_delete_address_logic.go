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

type UserDeleteAddressLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserDeleteAddressLogic(svcCtx *svc.ServiceContext) *UserDeleteAddressLogic {
	return &UserDeleteAddressLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteAddressLogic) UserDeleteAddress(ctx context.Context, req *types.IdPathReq) error {
	_, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/user/addresses/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, huser.NewAddressHandler(l.svcCtx).Delete)
	if err != nil {
		return err
	}
	return nil
}
