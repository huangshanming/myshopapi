package task

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type InternalEventLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalEventLogic {
	return &InternalEventLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalEventLogic) InternalEvent(ctx context.Context, req *types.TaskEventReq) (*types.EmptyResp, error) {
	if req.UserId == 0 {
		return nil, xerr.New(http.StatusBadRequest, "缺少 user_id")
	}
	code := strings.TrimSpace(req.Event)
	if code == "" {
		code = strings.TrimSpace(req.TaskCode)
	}
	if code == "" {
		return nil, xerr.New(http.StatusBadRequest, "缺少 task_code/event")
	}
	bizReq := biz.TaskEventReq{
		UserID:   req.UserId,
		TaskCode: code,
		Delta:    1,
		RefType:  strings.TrimSpace(req.RefType),
		RefID:    req.RefId.Uint64(),
	}
	if err := biz.NewTaskLogic(l.svcCtx).HandleEvent(ctx, bizReq); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
