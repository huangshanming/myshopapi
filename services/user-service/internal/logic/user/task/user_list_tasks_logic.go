package task

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListTasksLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListTasksLogic {
	return &UserListTasksLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserListTasksLogic) UserListTasks(ctx context.Context) (*types.PageListResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	list, err := biz.NewTaskLogic(l.svcCtx).ListUserTasks(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
