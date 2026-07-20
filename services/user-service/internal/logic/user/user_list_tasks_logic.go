package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
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

func (l *UserListTasksLogic) UserListTasks(w http.ResponseWriter, r *http.Request) {
	huser.NewTaskHandler(l.svcCtx).UserListTasks(w, r)
}
