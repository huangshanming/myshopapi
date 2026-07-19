package logic

import (
	"context"

	"mymall/services/inventory-sync-service/internal/svc"
	"mymall/services/inventory-sync-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthzLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthzLogic {
	return &HealthzLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *HealthzLogic) Healthz() (*types.HealthResp, error) {
	return &types.HealthResp{Status: "ok", Service: "inventory-sync-service"}, nil
}
