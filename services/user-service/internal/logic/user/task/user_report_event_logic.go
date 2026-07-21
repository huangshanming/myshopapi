package task

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type UserReportEventLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserReportEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserReportEventLogic {
	return &UserReportEventLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserReportEventLogic) UserReportEvent(ctx context.Context, req *types.TaskEventReq) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return xerr.New(http.StatusUnauthorized, "未登录")
	}
	refID, _ := strconv.ParseUint(req.RefId, 10, 64)
	bizReq := biz.TaskEventReq{
		UserID:   userID,
		TaskCode: req.Event,
		Delta:    1,
		RefID:    refID,
	}
	if err := biz.NewTaskLogic(l.svcCtx).HandleEvent(ctx, bizReq); err != nil {
		return xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil
}
