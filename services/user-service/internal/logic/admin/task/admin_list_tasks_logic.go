package task

import (
	"context"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListTasksLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListTasksLogic(svcCtx *svc.ServiceContext) *AdminListTasksLogic {
	return &AdminListTasksLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminListTasksLogic) AdminListTasks(ctx context.Context) (resp *types.PageListResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/tasks", nil, nil, nil, hadmin.NewTaskHandler(l.svcCtx).AdminList)
	if err != nil {
		return nil, err
	}
	var list interface{}
	if err := httpinvoke.Decode(raw, &list); err != nil {
		return nil, err
	}
	return &types.PageListResp{List: list}, nil
}
