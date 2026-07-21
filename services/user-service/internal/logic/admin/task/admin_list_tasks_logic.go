package task

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
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
	list, err := biz.NewTaskLogic(l.svcCtx).AdminListTasks(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
