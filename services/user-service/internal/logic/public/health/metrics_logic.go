package health

import (
	"context"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetricsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMetricsLogic(svcCtx *svc.ServiceContext) *MetricsLogic {
	return &MetricsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MetricsLogic) Metrics(ctx context.Context) error {
	return nil
}
