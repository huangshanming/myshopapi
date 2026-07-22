// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateTaskLogic {
	return &AdminUpdateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateTaskLogic) AdminUpdateTask(req *types.UpdateTaskReq) (resp *types.TaskResp, err error) {
	// todo: add your logic here and delete this line

	return
}
