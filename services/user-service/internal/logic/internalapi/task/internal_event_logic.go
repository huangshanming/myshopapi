package task

import (
	"context"
	"net/http"
	"strconv"

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

func (l *InternalEventLogic) InternalEvent(ctx context.Context, req *types.TaskEventReq) error {
	refID, _ := strconv.ParseUint(req.RefId, 10, 64)
	bizReq := biz.TaskEventReq{
		TaskCode: req.Event,
		Delta:    1,
		RefID:    refID,
	}
	if err := biz.NewTaskLogic(l.svcCtx).HandleEvent(ctx, bizReq); err != nil {
		return xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil
}
