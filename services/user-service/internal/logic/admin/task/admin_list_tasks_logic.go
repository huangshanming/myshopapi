package task

import (
	"context"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListTasksLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListTasksLogic {
	return &AdminListTasksLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListTasksLogic) AdminListTasks(ctx context.Context) (resp *types.PageListResp, err error) {
	data, err := hadmin.NewTaskHandler(l.svcCtx).AdminList(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.PageListResp{List: data}, nil
}
