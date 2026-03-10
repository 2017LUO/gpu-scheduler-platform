package logging

import (
	"go.uber.org/zap"
)

func WithComponent(lg *zap.Logger, component string) *zap.Logger {
	if lg == nil {
		return zap.NewNop()
	}
	return lg.With(zap.String("component", component))
}

func LoggerOrNop(lg *zap.Logger) *zap.Logger {
	if lg == nil {
		return zap.NewNop()
	}
	return lg
}
