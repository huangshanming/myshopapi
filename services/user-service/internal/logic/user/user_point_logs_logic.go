package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserPointLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPointLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPointLogsLogic {
	return &UserPointLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserPointLogsLogic) UserPointLogs(w http.ResponseWriter, r *http.Request) {
	huser.NewTaskHandler(l.svcCtx).UserPointLogs(w, r)
}
