package task

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListTasksLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserListTasksLogic(svcCtx *svc.ServiceContext) *UserListTasksLogic {
	return &UserListTasksLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserListTasksLogic) UserListTasks(ctx context.Context) (resp *types.PageListResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/user/tasks", nil, nil, nil, huser.NewTaskHandler(l.svcCtx).UserListTasks)
	if err != nil {
		return nil, err
	}
	var list interface{}
	if err := httpinvoke.Decode(raw, &list); err != nil {
		return nil, err
	}
	return &types.PageListResp{List: list}, nil
}
