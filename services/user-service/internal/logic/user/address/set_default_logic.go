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

type SetDefaultLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSetDefaultLogic(svcCtx *svc.ServiceContext) *SetDefaultLogic {
	return &SetDefaultLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *SetDefaultLogic) SetDefault(ctx context.Context, req *types.IdPathReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/user/addresses/{Id}/default", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, huser.NewAddressHandler(l.svcCtx).SetDefault)
	if err != nil {
		return err
	}
	return nil
}
