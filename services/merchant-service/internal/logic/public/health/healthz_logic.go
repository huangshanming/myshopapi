package health

import (
	"context"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthzLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewHealthzLogic(svcCtx *svc.ServiceContext) *HealthzLogic {
	return &HealthzLogic{Logger: logx.WithContext(context.Background()), svcCtx: svcCtx}
}

func (l *HealthzLogic) Healthz(ctx context.Context) (resp *types.EmptyResp, err error) {
	_ = ctx

	return &types.EmptyResp{}, nil
}
