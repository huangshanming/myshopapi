package task

import (
	"context"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListTasksLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListTasksLogic {
	return &UserListTasksLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserListTasksLogic) UserListTasks(ctx context.Context) (resp *types.PageListResp, err error) {
	data, err := huser.NewTaskHandler(l.svcCtx).UserListTasks(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.PageListResp{List: data}, nil
}
