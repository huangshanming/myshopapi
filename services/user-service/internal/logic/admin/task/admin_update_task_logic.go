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

type AdminUpdateTaskLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateTaskLogic {
	return &AdminUpdateTaskLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateTaskLogic) AdminUpdateTask(ctx context.Context, req *types.UpdateTaskReq) (resp *types.AnyResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "任务ID无效")
	}
	bizReq := biz.UpdateTaskReq{}
	if req.Name != "" {
		bizReq.Title = &req.Name
	}
	if req.Description != "" {
		bizReq.Description = &req.Description
	}
	if req.Points != 0 {
		p := int(req.Points)
		bizReq.RewardPoints = &p
	}
	if req.Status != "" {
		var enabled int8
		switch req.Status {
		case "1", "enabled", "true":
			enabled = 1
		default:
			enabled = 0
		}
		bizReq.Enabled = &enabled
	}
	t, err := biz.NewTaskLogic(l.svcCtx).AdminUpdateTask(ctx, req.Id, bizReq)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: t}, nil
}
