package logic

import (
	"net/http"

	"context"

	"mymall/pkg/metrics"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetricsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetricsLogic {
	return &MetricsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MetricsLogic) Metrics(w http.ResponseWriter, r *http.Request) {
	metrics.Handler()(w, r)
}
