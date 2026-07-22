// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListTasksLogic {
	return &UserListTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListTasksLogic) UserListTasks() (resp *types.PageListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
