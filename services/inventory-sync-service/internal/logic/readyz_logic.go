package logic

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/inventory-sync-service/internal/svc"
	"mymall/services/inventory-sync-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadyzLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadyzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadyzLogic {
	return &ReadyzLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ReadyzLogic) Readyz() (*types.ReadyResp, error) {
	if l.svcCtx.Health == nil {
		return &types.ReadyResp{Status: "ready"}, nil
	}
	ok, failed := l.svcCtx.Health.CheckAll(l.ctx)
	if !ok {
		return &types.ReadyResp{Status: "not_ready", Checks: failed}, xerr.New(http.StatusServiceUnavailable, "not_ready")
	}
	return &types.ReadyResp{Status: "ready"}, nil
}
