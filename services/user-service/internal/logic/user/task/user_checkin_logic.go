package task

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCheckinLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCheckinLogic(svcCtx *svc.ServiceContext) *UserCheckinLogic {
	return &UserCheckinLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserCheckinLogic) UserCheckin(ctx context.Context) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/tasks/checkin", nil, nil, nil, huser.NewTaskHandler(l.svcCtx).UserCheckin)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
