package task

import (
	"context"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCheckinLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCheckinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCheckinLogic {
	return &UserCheckinLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserCheckinLogic) UserCheckin(ctx context.Context) (resp *types.AnyResp, err error) {
	data, err := huser.NewTaskHandler(l.svcCtx).UserCheckin(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
