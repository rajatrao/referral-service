package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"infra",
	fx.Provide(
		NewLogger,
	),
)

func NewLogger() *zap.Logger {
	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.FunctionKey = "method"
	logger := zap.Must(logCfg.Build())

	return logger
}
