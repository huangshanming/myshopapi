package telemetry

import (
	"context"

	"mymall/pkg/config"
)

// Init 初始化 OpenTelemetry 导出器骨架；enabled=false 时为 no-op。
// 生产环境可替换为 otlptracehttp + Jaeger Collector。
func Init(ctx context.Context, cfg config.TelemetryConfig) (func(context.Context) error, error) {
	_ = ctx
	if !cfg.Enabled {
		return func(context.Context) error { return nil }, nil
	}
	// 骨架：启用时仅记录配置，后续可接入 go.opentelemetry.io/otel
	return func(context.Context) error { return nil }, nil
}
