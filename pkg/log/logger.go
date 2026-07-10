package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(service string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger.With(zap.String("service", service)), nil
}
