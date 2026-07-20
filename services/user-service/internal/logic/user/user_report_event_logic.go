package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserReportEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserReportEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserReportEventLogic {
	return &UserReportEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserReportEventLogic) UserReportEvent(w http.ResponseWriter, r *http.Request) {
	huser.NewTaskHandler(l.svcCtx).UserReportEvent(w, r)
}
