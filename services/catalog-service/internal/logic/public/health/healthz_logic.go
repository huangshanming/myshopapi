package health

import (
	"context"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthzLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewHealthzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthzLogic {
	return &HealthzLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *HealthzLogic) Healthz(ctx context.Context) (resp *types.EmptyResp, err error) {
	_ = ctx

	return &types.EmptyResp{}, nil
}
