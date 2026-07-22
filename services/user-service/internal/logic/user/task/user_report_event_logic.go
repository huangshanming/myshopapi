package task

import (
	"context"
	"net/http"
	"strings"

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

func (l *UserReportEventLogic) UserReportEvent(ctx context.Context, req *types.TaskEventReq) (*types.EmptyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	bizReq := toBizTaskEvent(req, userID)
	if err := biz.NewTaskLogic(l.svcCtx).HandleEvent(ctx, bizReq); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}

func toBizTaskEvent(req *types.TaskEventReq, userID uint64) biz.TaskEventReq {
	code := strings.TrimSpace(req.Event)
	if code == "" {
		code = strings.TrimSpace(req.TaskCode)
	}
	refID := req.RefId.Uint64()
	return biz.TaskEventReq{
		UserID:   userID,
		TaskCode: code,
		Delta:    1,
		RefType:  strings.TrimSpace(req.RefType),
		RefID:    refID,
	}
}
